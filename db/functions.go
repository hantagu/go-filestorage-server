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

func GetUserByPublicKey(public_key ed25519.PublicKey) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	result := &db_types.User{}
	if err := UsersCollection.FindOne(ctx, &bson.D{{Key: "public_key", Value: public_key}}).Decode(result); err != nil {
		return "", err
	}

	return result.Username, nil
}

func GetUserByUsername(username string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	result := &db_types.User{}
	if err := UsersCollection.FindOne(ctx, &bson.D{{Key: "username", Value: username}}).Decode(result); err != nil {
		return nil, err
	}

	return result.PublicKey, nil

}

func SetUsername(public_key ed25519.PublicKey, username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	update_result, err := UsersCollection.UpdateOne(ctx,

		&bson.D{{Key: "public_key", Value: public_key}},

		&bson.D{{Key: "$set", Value: &bson.D{
			{Key: "username", Value: username},
		}}},
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

func GetFileMetadata(public_key ed25519.PublicKey, name string) (*db_types.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	file_metadata := &db_types.File{}
	err := FilesCollection.FindOne(ctx, &bson.D{
		{Key: "$and", Value: &bson.A{
			&bson.D{{Key: "$or", Value: &bson.A{
				&bson.D{{Key: "owner", Value: public_key}},
				&bson.D{{Key: "access", Value: public_key}},
			}}},
			&bson.D{{Key: "name", Value: name}},
		}},
	}).Decode(file_metadata)

	if err != nil {
		return nil, err
	}

	return file_metadata, nil
}

func GetAllFilesMetadata(public_key ed25519.PublicKey) ([]db_types.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	cur, err := FilesCollection.Find(ctx, &bson.D{
		{Key: "$or", Value: &bson.A{
			&bson.D{{Key: "owner", Value: public_key}},
			&bson.D{{Key: "access", Value: public_key}},
		}},
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	result := []db_types.File{}
	err = cur.All(ctx, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
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

func DeleteFileMetadataByID(ID primitive.ObjectID) (*db_types.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	file_metadata := &db_types.File{}
	err := FilesCollection.FindOneAndDelete(ctx, &bson.D{{Key: "_id", Value: ID}}).Decode(file_metadata)

	if err != nil {
		return nil, err
	}

	return file_metadata, nil
}

func DeleteFileMetadata(file_owner ed25519.PublicKey, file_name string) (*db_types.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	file_metadata := &db_types.File{}
	err := FilesCollection.FindOneAndDelete(ctx,
		&bson.D{{Key: "$and", Value: &bson.A{
			&bson.D{{Key: "owner", Value: file_owner}},
			&bson.D{{Key: "name", Value: file_name}},
		}}},
	).Decode(file_metadata)

	if err != nil {
		return nil, err
	}

	return file_metadata, nil
}

func GrantAccess(file_owner ed25519.PublicKey, file_name string, public_key ed25519.PublicKey) (int, int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	update_result, err := FilesCollection.UpdateOne(ctx,

		&bson.D{{Key: "$and", Value: &bson.A{
			&bson.D{{Key: "owner", Value: file_owner}},
			&bson.D{{Key: "name", Value: file_name}},
		}}},

		&bson.D{{Key: "$addToSet", Value: &bson.D{
			{Key: "access", Value: public_key},
		}}},
	)

	return int(update_result.MatchedCount), int(update_result.ModifiedCount), err
}

func RevokeAccess(file_owner ed25519.PublicKey, file_name string, public_key ed25519.PublicKey) (int, int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()

	update_result, err := FilesCollection.UpdateOne(ctx,

		&bson.D{{Key: "$and", Value: &bson.A{
			&bson.D{{Key: "owner", Value: file_owner}},
			&bson.D{{Key: "name", Value: file_name}},
		}}},

		&bson.D{{Key: "$pull", Value: &bson.D{
			{Key: "access", Value: public_key},
		}}},
	)

	return int(update_result.MatchedCount), int(update_result.ModifiedCount), err
}
