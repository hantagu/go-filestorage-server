package handlers

import (
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

func ListFiles(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.Empty{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	result, err := db.GetAllFilesMetadata(request.PublicKey)

	if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	protocol.SendResponse(conn, true, &bson.D{{Key: "files", Value: result}})
}
