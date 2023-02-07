package mongodb

import (
	"context"
	"go-filestorage-server/utils"
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

func InitMongoDB() {

	// Create new client with URI from config
	opts := options.Client()
	opts.ApplyURI(utils.Config.MongoDB_URI)
	var err error
	if Client, err = mongo.NewClient(opts); err != nil {
		utils.Logger.Fatalln(err)
		return
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Connect(ctx); err != nil {
		utils.Logger.Fatalln(err)
		return
	}

	// Ping MongoDB server
	ctx, cancel = context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Ping(ctx, nil); err != nil {
		utils.Logger.Fatalln(err)
		return
	}

	Database = Client.Database(utils.Config.MongoDB_DB)
	UsersCollection = Database.Collection(utils.Config.MongoDB_UsersCollection)
	FilesCollection = Database.Collection(utils.Config.MongoDB_FilesCollection)

	// Make the `public_key` field unique in Users collection
	ctx, cancel = context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "public_key", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Make the `login` field unique in Users collection
	ctx, cancel = context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Make the `name` field unique for each user
	ctx, cancel = context.WithTimeout(context.Background(), utils.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	FilesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "owner", Value: 1}, {Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
