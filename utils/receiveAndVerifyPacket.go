package utils

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"go-filestorage-server/config"
	"go-filestorage-server/protocol"
	"io"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrPacketSignature = errors.New("invalid packet signature")

func ReceiveAndVerifyPacket(conn net.Conn) (*protocol.Request, error) {

	buffer := make([]byte, config.PROTO_BSON_DOCUMENT_LENGTH_SIZE)

	// Read first N bytes (according to BSON documentation) which indicate the size of the entire BSON document
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	// Convert bytes to a UInt32 (4-byte) value
	packetLength := binary.LittleEndian.Uint32(buffer) - config.PROTO_BSON_DOCUMENT_LENGTH_SIZE

	// Read BSON document
	buffer = append(buffer, make([]byte, packetLength)...)
	if _, err := io.ReadFull(conn, buffer[config.PROTO_BSON_DOCUMENT_LENGTH_SIZE:]); err != nil {
		return nil, err
	}

	// Decode a BSON document to a generic packet to find out it's type
	packet := protocol.Request{}
	if err := bson.Unmarshal(buffer, &packet); err != nil {
		return nil, err
	}

	// Receiving a signature bytes
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := io.ReadFull(conn, signature); err != nil {
		return nil, err
	}

	// Check a packet signature
	if !ed25519.Verify(packet.PublicKey, buffer, signature) {
		return nil, ErrPacketSignature
	}

	return &packet, nil
}
