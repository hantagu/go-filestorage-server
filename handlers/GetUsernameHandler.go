package handlers

import (
	"context"
	"errors"
	"go-filestorage-server/mongodb"
	"go-filestorage-server/mongodb/types"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleGetUsername(conn net.Conn, packet *protocol.Packet) {

	// Unmarshal packet data
	request := &protocol.GetUsername{}
	if err := bson.Unmarshal(packet.Data, request); err != nil {
		utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		protocol.SendDescriptionError(conn, "Failed to decode packet data")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	find_result := mongodb.UsersCollection.FindOne(ctx, nil)
	result := &types.User{}

	if err := find_result.Decode(result); errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendDescriptionError(conn, "This public key does not match any username")
		return
	} else if err != nil {
		protocol.SendDescriptionError(conn, "Internal error")
		return
	}

}
