package main

import (
	"bytes"
	"go-filestorage-server/handlers"
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
	preamble := make([]byte, utils.PROTO_NET_PREAMBLE_SIZE)
	if _, err := io.ReadFull(conn, preamble); err != nil {
		utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	} else if !bytes.Equal([]byte(utils.PROTO_NET_PREAMBLE), preamble) {
		utils.Logger.Printf("%s: wrong preamble, connection closed\n", conn.RemoteAddr())
		return
	}

	// Receive first packet in connection
	packet, err := utils.ReceiveAndVerifyPacket(conn)
	if err != nil {
		utils.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	}

	// Select handler function depending on the type of package
	switch packet.Type {
	case protocol.CLAIM_USERNAME:
		handlers.HandleClaimUsername(conn, packet)
	case protocol.UPLOAD_METADATA:
		handlers.HandleUpload(conn, packet)
	}
}
