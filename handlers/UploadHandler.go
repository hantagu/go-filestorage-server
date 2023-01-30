package handlers

import (
	"fmt"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleUpload(connection net.Conn, first_packet *protocol.Packet) {

	filepath := fmt.Sprintf("%s%c%d", utils.Config.UserdataPath, os.PathSeparator, time.Now().UnixNano())

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
}
