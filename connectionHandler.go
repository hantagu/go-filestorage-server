package main

import (
	"bytes"
	"go-filestorage-server/config"
	"go-filestorage-server/handlers"
	"go-filestorage-server/logger"
	"go-filestorage-server/protocol"
	"go-filestorage-server/utils"
	"io"
	"net"
	"sync"
)

func handleConnection(conn net.Conn, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	defer conn.Close()

	// Check a connection preamble
	preamble := make([]byte, config.PROTO_NET_PREAMBLE_SIZE)
	if _, err := io.ReadFull(conn, preamble); err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	} else if !bytes.Equal([]byte(config.PROTO_NET_PREAMBLE), preamble) {
		logger.Logger.Printf("%s: wrong preamble, connection closed\n", conn.RemoteAddr())
		return
	}

	// Receive first packet in connection
	packet, err := utils.ReceiveAndVerifyPacket(conn)
	if err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	}

	// Select handler function depending on the type of package
	switch packet.Type {
	case protocol.REQ_GET_USERNAME:
		handlers.GetUsername(conn, packet)
	case protocol.REQ_SET_USERNAME:
		handlers.SetUsername(conn, packet)
	default:
		protocol.SendDescription(conn, false, "Invalid request type")
	}
}
