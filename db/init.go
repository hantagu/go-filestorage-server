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

	// Создание нового клиента для подключения к MongoDB
	opts := options.Client()
	opts.ApplyURI(config.Config.MongoDB_URI)
	var err error
	if Client, err = mongo.NewClient(opts); err != nil {
		return err
	}

	// Подключение к серверу MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Connect(ctx); err != nil {
		return err
	}

	// Проверка соединения
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	if err := Client.Ping(ctx, nil); err != nil {
		return err
	}

	Database = Client.Database(config.Config.MongoDB_DB)
	UsersCollection = Database.Collection(config.Config.MongoDB_Users_Collection)
	FilesCollection = Database.Collection(config.Config.MongoDB_Files_Collection)

	// Создание индекса на поле `public_key`, который сделает его уникальным во всей коллекции пользователей
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "public_key", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Создание индекса на поле `username`, который сделает его уникальным во всей коллекции пользователей
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Создание индекса на поле `name`, который сделает его уникальным во всей коллекции файлов для каждого пользователя
	ctx, cancel = context.WithTimeout(context.Background(), config.MONGODB_CONTEXT_TIMEOUT*time.Second)
	defer cancel()
	FilesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "owner", Value: 1}, {Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return nil
}
