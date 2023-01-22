package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"go-filesharing-server/types"
	"io"
	"net"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrPacketSignature = errors.New("invalid signature")

func handleConnection(connection net.Conn, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	defer connection.Close()

	// Check a connection preamble
	preamble := make([]byte, PROTO_NET_PREAMBLE_SIZE)
	if _, err := io.ReadFull(connection, preamble); err != nil {
		logger.Printf("%s: %s\n", connection.RemoteAddr(), err.Error())
		return
	} else if !bytes.Equal([]byte(PROTO_NET_PREAMBLE), preamble) {
		logger.Printf("%s: Wrong preamble\n", connection.RemoteAddr())
		return
	}

	packet, err := receiveAndVerifyPacket(connection)
	if err != nil {
		logger.Printf("%s: %s\n", connection.RemoteAddr(), err.Error())
		return
	}

	switch packet.Type {
	// TODO
	}
}

func receiveAndVerifyPacket(connection net.Conn) (*types.Packet, error) {

	packetBytes := make([]byte, PROTO_BSON_DOCUMENT_LENGTH_SIZE)

	// Read first N bytes (according to BSON documentation) which indicate the size of the entire BSON document
	if _, err := io.ReadFull(connection, packetBytes); err != nil {
		logger.Printf("%s: %s\n", connection.RemoteAddr(), err.Error())
		return nil, err
	}

	// Convert bytes to a UInt32 (4-byte) value
	packetLength := binary.LittleEndian.Uint32(packetBytes)

	// Create a buffer and run loop to read the entire BSON document
	buffer := make([]byte, PROTO_NET_BUFFER_SIZE)
	for uint32(len(packetBytes)) < packetLength {
		if n, err := connection.Read(buffer); err != nil {
			logger.Printf("%s: %s\n", connection.RemoteAddr(), err.Error())
			return nil, err
		} else {
			packetBytes = append(packetBytes, buffer[:n]...)
		}
	}

	// Decode a BSON document to a generic packet to find out it's type
	packet := types.Packet{}
	if err := bson.Unmarshal(packetBytes, &packet); err != nil {
		logger.Printf("#54 %s: %s\n", connection.RemoteAddr(), err.Error())
		return nil, &net.DNSConfigError{}
	}

	// Receiving a signature bytes
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := io.ReadFull(connection, signature); err != nil {
		logger.Printf("%s: Failed to get signature\n", connection.RemoteAddr())
		return nil, err
	}

	// Check a packet signature
	if !ed25519.Verify(packet.PublicKey, packetBytes, signature) {
		logger.Printf("%s: Invalid packet signature\n", connection.RemoteAddr())
		return nil, ErrPacketSignature
	}

	return &packet, nil
}
