package utils

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"go-filestorage-server/protocol"
	"io"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrPacketSignature = errors.New("invalid packet signature")

func ReceiveAndVerifyPacket(connection net.Conn) (*protocol.Packet, error) {

	buffer := make([]byte, PROTO_BSON_DOCUMENT_LENGTH_SIZE)

	// Read first N bytes (according to BSON documentation) which indicate the size of the entire BSON document
	if _, err := io.ReadFull(connection, buffer); err != nil {
		return nil, err
	}

	// Convert bytes to a UInt32 (4-byte) value
	packetLength := binary.LittleEndian.Uint32(buffer) - PROTO_BSON_DOCUMENT_LENGTH_SIZE

	// Read BSON document
	buffer = append(buffer, make([]byte, packetLength)...)
	if _, err := io.ReadFull(connection, buffer[PROTO_BSON_DOCUMENT_LENGTH_SIZE:]); err != nil {
		return nil, err
	}

	// Decode a BSON document to a generic packet to find out it's type
	packet := protocol.Packet{}
	if err := bson.Unmarshal(buffer, &packet); err != nil {
		return nil, err
	}

	// Receiving a signature bytes
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := io.ReadFull(connection, signature); err != nil {
		return nil, err
	}

	// Check a packet signature
	if !ed25519.Verify(packet.PublicKey, buffer, signature) {
		return nil, ErrPacketSignature
	}

	return &packet, nil
}
