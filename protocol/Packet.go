package protocol

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
)

// Generic protocol packet
type Packet struct {
	Type      string            `bson:"type"`       // Type of packet
	PublicKey ed25519.PublicKey `bson:"public_key"` // Sender's public key
	Data      bson.Raw          `bson:"data"`       // Data of packet
}
