package handlers

import (
	"bytes"
	"crypto/ed25519"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

func GrantAccess(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.FileAccess{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	if len(request_data.PublicKey) != ed25519.PublicKeySize {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid size of public key"})
		return
	} else if bytes.Equal(request.PublicKey, request_data.PublicKey) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "You cannot grant access to your file to yourself"})
		return
	}

	// Insert the public key into the file's metadata in the database
	if matched, modified, err := db.GrantAccess(request.PublicKey, request_data.Name, request_data.PublicKey); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	} else if matched == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "A file with this name does not exist"})
	} else if modified == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "This user already has access to this file"})
	} else {
		protocol.SendResponse(conn, true, &protocol.Empty{})
	}
}
