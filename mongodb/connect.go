package mongodb

import (
	"context"
	"go-filestorage-server/utils"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDB() {

	opts := options.Client()
	opts.ApplyURI(utils.Config.MongoDB_URI)

	var err error

	if Client, err = mongo.NewClient(opts); err != nil {
		utils.Logger.Fatalln(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Client.Connect(ctx); err != nil {
		utils.Logger.Fatalln(err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, nil); err != nil {
		utils.Logger.Fatalln(err)
		return
	}
}
