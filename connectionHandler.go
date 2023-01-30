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

func handleConnection(connection net.Conn, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	defer connection.Close()

	// Check a connection preamble
	preamble := make([]byte, utils.PROTO_NET_PREAMBLE_SIZE)
	if _, err := io.ReadFull(connection, preamble); err != nil {
		utils.Logger.Printf("%s: %s\n", connection.RemoteAddr(), err)
		return
	} else if !bytes.Equal([]byte(utils.PROTO_NET_PREAMBLE), preamble) {
		utils.Logger.Printf("%s: wrong preamble, connection closed\n", connection.RemoteAddr())
		return
	}

	packet, err := utils.ReceiveAndVerifyPacket(connection)
	if err != nil {
		utils.Logger.Printf("%s: %s\n", connection.RemoteAddr(), err)
		return
	}

	switch packet.Type {
	case protocol.UPLOAD_METADATA:
		handlers.HandleUpload(connection, packet)
	}
}
