package handlers

import (
	"errors"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsername(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.PublicKey{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Find the username in the database
	if result, err := db.GetUsername(request.PublicKey); errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "This public key does not match any username"})
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	} else {
		protocol.SendResponse(conn, true, &protocol.Username{Username: result})
	}
}
