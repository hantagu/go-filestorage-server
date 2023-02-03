package handlers

import (
	"context"
	"go-filestorage-server/mongodb"
	"go-filestorage-server/mongodb/types"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleClaimUsername(connection net.Conn, packet *protocol.Packet) {

	request := &protocol.ClaimUsername{}
	if err := bson.Unmarshal(packet.Data, request); err != nil {
		utils.Logger.Printf("%s: %s\n", connection.RemoteAddr(), err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	if _, err := mongodb.UsersCollection.InsertOne(ctx, types.User{PublicKey: packet.PublicKey, Username: request.Username}); mongo.IsDuplicateKeyError(err) {

		description, _ := bson.Marshal(protocol.ResponseDescription{Description: "This username is already taken"})
		resp, _ := bson.Marshal(protocol.Response{Successful: false, Data: description})

		connection.Write(resp)
	}
}
