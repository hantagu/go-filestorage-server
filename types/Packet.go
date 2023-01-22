package types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
)

// Generic protcol packet
type Packet struct {
	Type      string            `bson:"type"`       // Type of packet
	PublicKey ed25519.PublicKey `bson:"public_key"` // Sender's public key
	Head      bson.Raw          `bson:"head"`       // Packet head
	Body      bson.Raw          `bson:"body"`       // Packet body
}
