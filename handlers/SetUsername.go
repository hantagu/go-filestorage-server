package handlers

import (
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUsername(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.SetUsernameRequestData{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
		return
	}

	// Check username with regular expression
	if !regexp.MustCompile(`[a-z0-9_]{5,}`).MatchString(request_data.Username) {
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: "You can only use the characters a-z, 0-9 and the underscore character and minimum length is 5 characters"})
		return
	}

	if err := db.SetUsername(request.PublicKey, request_data.Username); mongo.IsDuplicateKeyError(err) {
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: "This username is already taken"})
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.ResponseDescription{Description: err.Error()})
	} else {
		protocol.SendResponse(conn, true, &protocol.ResponseDescription{Description: "Username changed successfully"})
	}
}
