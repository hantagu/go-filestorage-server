package handlers

import (
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

func GrantAccess(conn net.Conn, request *protocol.Request) {

	request_data := &protocol.PublicKey{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}
}
