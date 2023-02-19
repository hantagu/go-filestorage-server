package handlers

import (
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/logger"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadFile(conn net.Conn, upload_metadata_request *protocol.Request) {

	file_id := primitive.NewObjectID()
	file_path := fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_id.Hex())

	file, err := os.Create(file_path)
	if err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
		return
	}
	defer file.Close()

	// Unmarshal file metadata
	upload_metadata_data := &protocol.UploadFileMetadataRequestData{}
	if err := bson.Unmarshal(upload_metadata_request.Data, upload_metadata_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
		return
	}

	// Try to insert file metadata to MongoDB
	if err := db.InsertFileMetadata(file_id, upload_metadata_request.PublicKey, upload_metadata_data.Name); err != nil {

		file.Close()
		os.Remove(file_path)
		db.DeleteFileMetadata(upload_metadata_request.PublicKey, upload_metadata_data.Name)

		if mongo.IsDuplicateKeyError(err) {
			protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: "A file with this name already exists"})
			return
		}

		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
		return
	}

	// If the metadata is successfully inserted into MongoDB, then answer that the server is ready to receive the file
	protocol.SendResponse(conn, false, &protocol.UploadFileMetadataResponseData{})

	// Read all parts of a file
	var i uint32 = 0
	for ; i < upload_metadata_data.Parts; i++ {

		upload_content_request, err := utils.ReceiveAndVerifyPacket(conn)

		if err != nil || upload_content_request.Type != protocol.REQ_UPLOAD_FILE {
			file.Close()
			os.Remove(file_path)
			db.DeleteFileMetadata(upload_metadata_request.PublicKey, upload_metadata_data.Name)
			protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
			break
		}

		upload_content_data := &protocol.UploadFileContentRequestData{}
		if err := bson.Unmarshal(upload_content_request.Data, upload_content_data); err != nil {
			file.Close()
			os.Remove(file_path)
			db.DeleteFileMetadata(upload_metadata_request.PublicKey, upload_metadata_data.Name)
			protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
			break
		}

		file.Write(upload_content_data.Content)
		protocol.SendResponse(conn, true, &protocol.UploadFileContentResponseData{})
	}
}
