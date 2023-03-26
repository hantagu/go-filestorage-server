package db_types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`        // Unique user ID
	PublicKey ed25519.PublicKey  `bson:"public_key"` // User's public key
	Username  string             `bson:"username"`   // Unique username associated with public key
}
