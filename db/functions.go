package db

import (
	"context"
	"crypto/ed25519"
	"go-filestorage-server/config"
	db_types "go-filestorage-server/db/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsername(public_key ed25519.PublicKey) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	result := &db_types.User{}
	if err := UsersCollection.FindOne(ctx, bson.D{{Key: "public_key", Value: public_key}}).Decode(result); err != nil {
		return "", err
	}

	return result.Username, nil
}

func SetUsername(public_key ed25519.PublicKey, username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	update_result, err := UsersCollection.UpdateOne(ctx,
		bson.D{
			{Key: "public_key", Value: public_key},
		},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "username", Value: username},
			}},
		},
	)

	if err != nil {
		return err
	} else if update_result.ModifiedCount > 0 {
		return nil
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	_, err = UsersCollection.InsertOne(ctx,
		&db_types.User{
			ID:        primitive.NewObjectID(),
			PublicKey: public_key,
			Username:  username,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func GetFileMetadata(owner ed25519.PublicKey, name string) (*db_types.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	file_metadata := &db_types.File{}
	err := FilesCollection.FindOne(ctx, bson.D{
		{Key: "owner", Value: owner},
		{Key: "name", Value: name},
	}).Decode(file_metadata)

	if err != nil {
		return nil, err
	}

	return file_metadata, nil
}

func InsertFileMetadata(ID primitive.ObjectID, owner ed25519.PublicKey, name string, encrypted bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	_, err := FilesCollection.InsertOne(ctx, &db_types.File{
		ID:        ID,
		Owner:     owner,
		Name:      name,
		Encrypted: encrypted,
		Access:    []ed25519.PublicKey{},
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteFileMetadataByID(ID primitive.ObjectID) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	delete_result, err := FilesCollection.DeleteOne(ctx,
		bson.D{
			{Key: "_id", Value: ID},
		},
	)

	return int(delete_result.DeletedCount), err
}

func DeleteFileMetadata(file_owner ed25519.PublicKey, file_name string) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	delete_result, err := FilesCollection.DeleteOne(ctx,
		bson.D{
			{Key: "owner", Value: file_owner},
			{Key: "name", Value: file_name},
		},
	)

	return int(delete_result.DeletedCount), err
}
