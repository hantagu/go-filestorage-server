package main

import (
	"bytes"
	"errors"
	"go-filestorage-server/config"
	"go-filestorage-server/handlers"
	"go-filestorage-server/protocol"
	"io"
	"log"
	"net"
	"sync"
)

func handleConnection(conn net.Conn, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	defer conn.Close()

	// Check the connection preamble
	preamble := make([]byte, config.PROTOCOL_PREAMBLE_SIZE)
	if _, err := io.ReadFull(conn, preamble); err != nil {
		log.Default().Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	} else if !bytes.Equal([]byte(config.PROTOCOL_PREAMBLE), preamble) {
		log.Default().Printf("%s: wrong preamble, connection closed\n", conn.RemoteAddr())
		return
	}

	// Receive first request in this connection
	request, err := protocol.ReceiveAndVerifyPacket(conn)
	if errors.Is(err, protocol.ErrPacketSignature) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid request signature"})
	} else if err != nil {
		log.Default().Printf("%s: %s\n", conn.RemoteAddr(), err)
		return
	}

	// Select a handler function depending on the type of packet
	switch request.Type {
	case protocol.REQ_GET_USERNAME:
		handlers.GetUsername(conn, request)
	case protocol.REQ_SET_USERNAME:
		handlers.SetUsername(conn, request)
	case protocol.REQ_LIST_FILES:
		handlers.ListFiles(conn, request)
	case protocol.REQ_UPLOAD_FILE:
		handlers.UploadFile(conn, request)
	case protocol.REQ_DOWNLOAD_FILE:
		handlers.DownloadFile(conn, request)
    case protocol.REQ_DELETE_FILE:
		handlers.DeleteFile(conn, request)
	case protocol.REQ_GRANT_ACCESS:
		handlers.GrantAccess(conn, request)
	case protocol.REQ_REVOKE_ACCESS:
		handlers.RevokeAccess(conn, request)
	default:
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Invalid request type"})
	}
}
