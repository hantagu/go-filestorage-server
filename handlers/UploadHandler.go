package handlers

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"go-filestorage-server/mongodb"
	"go-filestorage-server/mongodb/types"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleUpload(connection net.Conn, first_packet *protocol.Packet) {

	file_id := primitive.NewObjectID()
	filepath := fmt.Sprintf("%s%c%s", utils.Config.UserdataPath, os.PathSeparator, file_id.Hex())

	file, err := os.Create(filepath)
	if err != nil {
		utils.Logger.Printf("%s: failed to create file in `%s` directory\n", connection.RemoteAddr(), utils.Config.UserdataPath)
		return
	}
	defer file.Close()

	// Unmarshal file metadata
	uploadMetadata := &protocol.UploadMetadata{}
	bson.Unmarshal(first_packet.Data, uploadMetadata)

	// Read all parts of a file
	var i uint32 = 0
	for ; i < uploadMetadata.Parts; i++ {
		packet, err := utils.ReceiveAndVerifyPacket(connection)
		if err != nil {
			utils.Logger.Printf("%s: %s\n", connection.RemoteAddr(), err)
			file.Close()
			os.Remove(filepath)
			return
		}

		uploadData := &protocol.UploadData{}
		bson.Unmarshal(packet.Data, uploadData)

		file.Write(uploadData.Content)
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	mongodb.FilesCollection.InsertOne(ctx, types.File{
		ID:     file_id,
		Owner:  first_packet.PublicKey,
		Name:   uploadMetadata.Name,
		Access: []ed25519.PublicKey{},
	})
}
