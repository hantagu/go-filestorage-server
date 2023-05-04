package db_types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`        // Уникальный ID пользователя, генерируемый MongoDB
	PublicKey ed25519.PublicKey  `bson:"public_key"` // Публичный ключ
	Username  string             `bson:"username"`   // Имя пользователя
}
