package protocol

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"go-filestorage-server/config"
	"io"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrPacketSignature = errors.New("invalid packet signature")

func ReceiveAndVerifyPacket(conn net.Conn) (*Request, error) {

	buffer := make([]byte, config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE)

	// Read first 4 bytes (according to BSON documentation) which indicate the size of the entire BSON document
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	// Convert bytes to a UInt32 (4-byte) value
	packetLength := binary.LittleEndian.Uint32(buffer) - config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE

	// Read BSON document
	buffer = append(buffer, make([]byte, packetLength)...)
	if _, err := io.ReadFull(conn, buffer[config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE:]); err != nil {
		return nil, err
	}

	// Decode a BSON document to a generic packet to find out it's type
	packet := Request{}
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

func SendResponse(conn net.Conn, successful bool, response_data interface{}) {

	raw_response_data, _ := bson.Marshal(response_data)

	raw_response, _ := bson.Marshal(&Response{
		Successful: successful,
		Data:       raw_response_data,
	})

	conn.Write(raw_response)
}
