package db

import (
	"context"
	"go-filestorage-server/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client          *mongo.Client
	Database        *mongo.Database
	UsersCollection *mongo.Collection
	FilesCollection *mongo.Collection
)

func Init() error {

	// Create a new client with URI from the config
	opts := options.Client()
	opts.ApplyURI(config.Config.MongoDB_URI)
	var err error
	if Client, err = mongo.NewClient(opts); err != nil {
		return err
	}

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Connect(ctx); err != nil {
		return err
	}

	// Ping the MongoDB server
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Ping(ctx, nil); err != nil {
		return err
	}

	Database = Client.Database(config.Config.MongoDB_DB)
	UsersCollection = Database.Collection(config.Config.MongoDB_UsersCollection)
	FilesCollection = Database.Collection(config.Config.MongoDB_FilesCollection)

	// Make the `public_key` field unique in the Users collection
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "public_key", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Make the `username` field unique in the Users collection
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Make the `name` field unique for each user
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	FilesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "owner", Value: 1}, {Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return nil
}
