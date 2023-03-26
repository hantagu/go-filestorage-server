package db_types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID  `bson:"_id"`       // Unique file ID
	Owner     ed25519.PublicKey   `bson:"owner"`     // Owner's public key
	Name      string              `bson:"name"`      // Name of file given by user
	Encrypted bool                `bson:"encrypted"` // Was the file additionally encrypted on the client side?
	Access    []ed25519.PublicKey `bson:"access"`    // Public keys with read access for this file
}
