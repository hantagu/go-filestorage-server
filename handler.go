package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"go-filesharing-server/types"
	"io"
	"net"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

func handleConnection(connection net.Conn, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	defer connection.Close()

	// Read a protocol preamble
	preamble := make([]byte, PROTO_NET_PREAMBLE_SIZE)
	if _, err := io.ReadFull(connection, preamble); err != nil {
		logger.Printf("#23 %s: %s\n", connection.RemoteAddr(), err.Error())
		return
	} else if !bytes.Equal([]byte(PROTO_NET_PREAMBLE), preamble) {
		logger.Printf("%s: Wrong preamble\n", connection.RemoteAddr())
		return
	}

	// Read first N bytes (according to BSON documentation) which indicate the size of the entire BSON document
	documentBytes := make([]byte, PROTO_BSON_DOCUMENT_LENGTH_SIZE)

	if _, err := io.ReadFull(connection, documentBytes); err != nil {
		logger.Printf("#34 %s: %s\n", connection.RemoteAddr(), err.Error())
		return
	}
	// Convert bytes to a UInt32 (4-byte) value
	documentLength := binary.LittleEndian.Uint32(documentBytes)

	// Create a buffer and run loop to read the entire BSON document
	buffer := make([]byte, PROTO_NET_BUFFER_SIZE)
	for uint32(len(documentBytes)) < documentLength {
		if n, err := connection.Read(buffer); err == nil {
			documentBytes = append(documentBytes, buffer[:n]...)
		} else {
			logger.Printf("#46 %s: %s\n", connection.RemoteAddr(), err.Error())
			return
		}
	}

	// Receiving a signature bytes
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := io.ReadFull(connection, signature); err != nil {
		logger.Printf("%s: Failed to get signature\n", connection.RemoteAddr())
		return
	}

	// Decode a BSON document to a generic packet to find out it's type
	document := types.Packet{}
	if err := bson.Unmarshal(documentBytes, &document); err != nil {
		logger.Printf("#54 %s: %s\n", connection.RemoteAddr(), err.Error())
		return
	}

	// Signature verification
	if !ed25519.Verify(document.PublicKey, documentBytes, signature) {
		logger.Printf("%s: Invalid packet signature\n", connection.RemoteAddr())
		return
	}
}
