package handlers

import (
	"errors"
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteFile(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.FileName{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Delete file's metadata from the database
	file_metadata, err := db.DeleteFileMetadata(request.PublicKey, request_data.Name)

	if errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "A file with this name does not exist"})
		return
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Delete a file
	os.Remove(fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_metadata.ID.Hex()))

	// Send a response
	protocol.SendResponse(conn, true, &protocol.Empty{})
}
