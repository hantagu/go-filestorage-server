package protocol

import (
	"crypto/ed25519"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

// Generic protocol incoming packet
type Packet struct {
	Type      string            `bson:"type"`       // Type of packet
	PublicKey ed25519.PublicKey `bson:"public_key"` // Sender's public key
	Data      bson.Raw          `bson:"data"`       // Data of packet
}

// Generic protocol response packet
type Response struct {
	Successful bool     `bson:"successful"` //
	Data       bson.Raw `bson:"data"`       //
}

// Response data with simple string description
type ResponseDescription struct {
	Description string `bson:"description"` // Description
}

func SendDescriptionError(conn net.Conn, desc string) {

	description := ResponseDescription{desc}
	raw_description, _ := bson.Marshal(description)

	response := Response{false, raw_description}
	raw_response, _ := bson.Marshal(response)

	conn.Write(raw_response)
}