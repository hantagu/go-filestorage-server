package handlers

import (
	"context"
	"errors"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	db_types "go-filestorage-server/db/types"
	"go-filestorage-server/logger"
	"go-filestorage-server/protocol"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsername(conn net.Conn, packet *protocol.Request) {

	// Unmarshal packet data
	request := &protocol.GetUsernameRequestData{}
	if err := bson.Unmarshal(packet.Data, request); err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		protocol.SendDescription(conn, false, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	result := &db_types.User{}

	if err := db.UsersCollection.FindOne(ctx, bson.D{{Key: "public_key", Value: request.PublicKey}}).Decode(result); errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendDescription(conn, false, "This public key does not match any username")
		return
	} else if err != nil {
		protocol.SendDescription(conn, false, err.Error())
		return
	}

	response_data, _ := bson.Marshal(&protocol.GetUsernameResponseData{Username: result.Username})
	response := &protocol.Response{Successful: true, Data: response_data}
	raw_response, _ := bson.Marshal(response)
	conn.Write(raw_response)
}
