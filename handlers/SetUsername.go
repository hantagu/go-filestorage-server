package handlers

import (
	"context"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	db_types "go-filestorage-server/db/types"
	"go-filestorage-server/protocol"
	"net"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUsername(conn net.Conn, request *protocol.Request) {

	// Unmarshal packet data
	request_data := &protocol.SetUsernameRequestData{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendDescription(conn, false, err.Error())
		return
	}

	if !regexp.MustCompile(`[a-z0-9_]{5,}`).MatchString(request_data.Username) {
		protocol.SendDescription(conn, false, "You can only use the characters a-z, 0-9 and the underscore character and minimum length is 5 characters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	update_result, err := db.UsersCollection.UpdateOne(ctx,
		// Filter
		bson.D{{
			Key:   "public_key",
			Value: request.PublicKey,
		}},
		// Update
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "username",
				Value: request_data.Username,
			}},
		}},
	)

	if err != nil {
		protocol.SendDescription(conn, false, err.Error())
		return
	}

	if update_result.ModifiedCount > 0 {
		protocol.SendDescription(conn, true, "Username set successfully")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	_, err = db.UsersCollection.InsertOne(ctx,
		&db_types.User{
			ID:        primitive.NewObjectID(),
			PublicKey: request.PublicKey,
			Username:  request_data.Username,
		},
	)

	if mongo.IsDuplicateKeyError(err) {
		protocol.SendDescription(conn, false, "This username is already taken")
	} else if err != nil {
		protocol.SendDescription(conn, false, err.Error())
	} else {
		protocol.SendDescription(conn, true, "Username set successfully")
	}
}
