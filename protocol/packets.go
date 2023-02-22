package protocol

import (
	"crypto/ed25519"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

func SendResponse(conn net.Conn, successful bool, response_data interface{}) {

	raw_response_data, _ := bson.Marshal(response_data)

	raw_response, _ := bson.Marshal(&Response{
		Successful: successful,
		Data:       raw_response_data,
	})

	conn.Write(raw_response)
}

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
type Empty struct {
}

// Packet data with simple string description
type Description struct {
	Description string `bson:"description"`
}

// -------------------------------------------------- Users -------------------------------------------------- //

// Packet data with public key
type PublicKey struct {
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Packet data with username
type Username struct {
	Username string `bson:"username"`
}

// -------------------------------------------------- Files -------------------------------------------------- //

// Packet data with file name
type FileName struct {
	Name string `bson:"name"`
}

// Packet data with file part
type FilePart struct {
	Part    uint32 `bson:"part"`
	Content []byte `bson:"content"`
}

// Packet data with file metadata
type FileMetadata struct {
	Name      string `bson:"name"`      // File name
	Encrypted bool   `bson:"encrypted"` // Was the file additionally encrypted on the client side?
	Parts     uint32 `bson:"parts"`     // Total parts
}
