package handlers

import (
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"os"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadFile(conn net.Conn, request *protocol.Request) {

	file_id := primitive.NewObjectID()
	file_path := fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_id.Hex())

	file, err := os.Create(file_path)
	if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}
	defer file.Close()

	// Unmarshal file metadata
	request_data := &protocol.FileMetadata{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		goto exit
	} else if request_data.Parts == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "File cannot be empty"})
		goto exit
	} else if !regexp.MustCompile(`^[^[:cntrl:]/]+$`).MatchString(request_data.Name) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid file name"})
		goto exit
	}

	// Try to insert file metadata to MongoDB
	if err := db.InsertFileMetadata(file_id, request.PublicKey, request_data.Name, request_data.Encrypted); err != nil {

		if mongo.IsDuplicateKeyError(err) {
			protocol.SendResponse(conn, false, &protocol.Description{Description: "A file with this name already exists"})
		} else {
			protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		}

		goto exit
	}

	// If the metadata is successfully inserted into MongoDB, then answer that the server is ready to receive the file
	protocol.SendResponse(conn, true, &protocol.Empty{})

	// Read all parts of a file
	for i := uint32(0); i < request_data.Parts; i++ {

		request, err := protocol.ReceiveAndVerifyPacket(conn)

		if err != nil || request.Type != protocol.REQ_UPLOAD_FILE {

			if err != nil {
				protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
			} else {
				protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid request"})
			}

			goto exit
		}

		request_data := &protocol.FilePart{}
		if err := bson.Unmarshal(request.Data, request_data); err != nil {
			protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
			goto exit
		}

		file.Write(request_data.Content)
		protocol.SendResponse(conn, true, &protocol.Empty{})
	}

	return

exit:
	os.Remove(file_path)
	db.DeleteFileMetadataByID(file_id)
}
