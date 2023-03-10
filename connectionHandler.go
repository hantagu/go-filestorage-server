package main

import (
	"bytes"
	"errors"
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
	preamble := make([]byte, config.PROTOCOL_PREAMBLE_SIZE)
	if _, err := io.ReadFull(conn, preamble); err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	} else if !bytes.Equal([]byte(config.PROTOCOL_PREAMBLE), preamble) {
		logger.Logger.Printf("%s: wrong preamble, connection closed\n", conn.RemoteAddr())
		return
	}

	// Receive first request in connection
	request, err := utils.ReceiveAndVerifyPacket(conn)
	if errors.Is(err, utils.ErrPacketSignature) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid request signature"})
	} else if err != nil {
		logger.Logger.Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	}

	// Select handler function depending on the type of package
	switch request.Type {
	case protocol.REQ_GET_USERNAME:
		handlers.GetUsername(conn, request)
	case protocol.REQ_SET_USERNAME:
		handlers.SetUsername(conn, request)
	case protocol.REQ_UPLOAD_FILE:
		handlers.UploadFile(conn, request)
	case protocol.REQ_DOWNLOAD_FILE:
		handlers.DownloadFile(conn, request)
	default:
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid request type"})
	}
}
