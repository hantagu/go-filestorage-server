package types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID     primitive.ObjectID  `bson:"_id"`
	Owner  ed25519.PublicKey   `bson:"owner"`
	Name   string              `bson:"name"`
	Access []ed25519.PublicKey `bson:"access"`
}
