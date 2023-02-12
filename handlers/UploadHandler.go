package handlers

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"go-filestorage-server/db"
	db_types "go-filestorage-server/db/types"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleUpload(conn net.Conn, first_packet *protocol.Packet) {

	file_id := primitive.NewObjectID()
	filepath := fmt.Sprintf("%s%c%s", utils.Config.UserdataPath, os.PathSeparator, file_id.Hex())

	file, err := os.Create(filepath)
	if err != nil {
		utils.Logger.Printf("%s: failed to create file in `%s` directory\n", conn.RemoteAddr(), utils.Config.UserdataPath)
		protocol.SendDescriptionError(conn, "Internal error")
		return
	}
	defer file.Close()

	// Unmarshal file metadata
	uploadMetadata := &protocol.UploadMetadata{}
	if err := bson.Unmarshal(first_packet.Data, uploadMetadata); err != nil {
		protocol.SendDescriptionError(conn, "Failed to decode packet data")
		return
	}

	// Read all parts of a file
	var i uint32 = 0
	for ; i < uploadMetadata.Parts; i++ {
		packet, err := utils.ReceiveAndVerifyPacket(conn)
		if err != nil {
			utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
			file.Close()
			os.Remove(filepath)
			protocol.SendDescriptionError(conn, "Internal error")
			return
		}

		uploadData := &protocol.UploadData{}
		if err := bson.Unmarshal(packet.Data, uploadData); err != nil {
			protocol.SendDescriptionError(conn, "Failed to decode packet data")
			return
		}

		file.Write(uploadData.Content)
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	if _, err := db.FilesCollection.InsertOne(ctx, db_types.File{
		ID:     file_id,
		Owner:  first_packet.PublicKey,
		Name:   uploadMetadata.Name,
		Access: []ed25519.PublicKey{},
	}); err != nil {
		utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		file.Close()
		os.Remove(filepath)
		protocol.SendDescriptionError(conn, "Internal error")
		return
	}
}
