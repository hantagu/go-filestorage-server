package handlers

import (
	"fmt"
	"go-filestorage-server/types"
	"go-filestorage-server/utils"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleUpload(connection net.Conn, first_packet *types.Packet) {

	filepath := fmt.Sprintf("%s%c%d", utils.Config.UserdataPath, os.PathSeparator, time.Now().UnixNano())

	file, err := os.Create(filepath)
	if err != nil {
		utils.Logger.Printf("%s: failed to create file in `%s` directory\n", connection.RemoteAddr(), utils.Config.UserdataPath)
		return
	}
	defer file.Close()

	uploadMetadata := &types.UploadMetadata{}
	bson.Unmarshal(first_packet.Data, uploadMetadata)

	utils.Logger.Printf("%+#v\n", uploadMetadata)

	var i uint32 = 0
	for ; i < uploadMetadata.Parts; i++ {
		packet, err := utils.ReceiveAndVerifyPacket(connection)
		if err != nil {
			utils.Logger.Printf("UH36 %s: %s\n", connection.RemoteAddr(), err)
			file.Close()
			os.Remove(filepath)
			return
		}

		uploadData := &types.UploadData{}
		bson.Unmarshal(packet.Data, uploadData)

		file.Write(uploadData.Content)
	}
}
