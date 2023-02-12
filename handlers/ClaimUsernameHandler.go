package handlers

import (
	"context"
	"go-filestorage-server/db"
	db_types "go-filestorage-server/db/types"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleClaimUsername(conn net.Conn, packet *protocol.Packet) {

	// Unmarshal packet data
	request := &protocol.ClaimUsername{}
	if err := bson.Unmarshal(packet.Data, request); err != nil {
		utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		protocol.SendDescriptionError(conn, "Failed to decode packet data")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	// Try to insert the information into database and send a response if there is an error (username is already taken)
	if _, err := db.UsersCollection.InsertOne(ctx, db_types.User{PublicKey: packet.PublicKey, Username: request.Username}); mongo.IsDuplicateKeyError(err) {
		protocol.SendDescriptionError(conn, "This username is already taken or you already have a username")
	} else if err != nil {
		protocol.SendDescriptionError(conn, "Internal error")
	}

	raw_response, err := bson.Marshal(protocol.Response{
		Successful: true,
		Data:       nil,
	})

	if err != nil {
		protocol.SendDescriptionError(conn, "Internal error")
		return
	}

	conn.Write(raw_response)
}
