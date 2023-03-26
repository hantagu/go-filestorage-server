package protocol

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
)

// Generic request packet
type Request struct {
	Type      string            `bson:"type"`       // Type of request
	PublicKey ed25519.PublicKey `bson:"public_key"` // Sender's public key
	Data      bson.Raw          `bson:"data"`       // Packet data
}

// Generic response packet
type Response struct {
	Successful bool     `bson:"successful"` // Was the request successful?
	Data       bson.Raw `bson:"data"`       // Packet data
}

// Empty packet data
type Empty struct{}

// Packet data with simple string description
type Description struct {
	Description string `bson:"description"`
}

// -------------------------------------------------- Users -------------------------------------------------- //

// Packet data with a public key
type PublicKey struct {
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Packet data with a username
type Username struct {
	Username string `bson:"username"`
}

// -------------------------------------------------- Files -------------------------------------------------- //

// Packet data with a file name
type FileName struct {
	Name string `bson:"name"`
}

// Packet data with file name and public key to which access to this file will be granted / revoked
type FileAccess struct {
	Name      string            `bson:"name"`
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Packet data with a chunk of the file
type FileChunk struct {
	Chunk   uint32 `bson:"chunk"`   // Number of this chunk
	Content []byte `bson:"content"` // Content of this chunk
}

// Packet data with a metadata of the file
type FileMetadata struct {
	Name      string `bson:"name"`      // File name
	Encrypted bool   `bson:"encrypted"` // Was this file additionally encrypted on the client side?
	Chunks    uint32 `bson:"chunks"`    // Total number of chunks
}
