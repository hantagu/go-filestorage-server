package handlers

import (
	"errors"
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"io"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DownloadFile(conn net.Conn, request *protocol.Request) {

	// Unmarshal request data
	request_data := &protocol.FileName{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Search file metadata in database
	file_metadata, err := db.GetFileMetadata(request.PublicKey, request_data.Name)
	if errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "A file with this name does not exist"})
		return
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Open a file
	file, err := os.Open(fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_metadata.ID.Hex()))
	if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Calculate a number of chunks
	if stat, err := file.Stat(); err == nil {
		parts := stat.Size() / config.PROTOCOL_CHUNK_SIZE
		if stat.Size()%config.PROTOCOL_CHUNK_SIZE != 0 {
			parts += 1
		}
		protocol.SendResponse(conn, true, &protocol.FileMetadata{Name: file_metadata.Name, Encrypted: false, Parts: uint32(parts)})
	} else {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	}

	buffer := make([]byte, 8*1024*1024)
	var i uint32 = 0
	for ; ; i++ {

		n, err := io.ReadFull(file, buffer)

		if n > 0 {
			protocol.SendResponse(conn, true, &protocol.FilePart{Part: i, Content: buffer[:n]})
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

}
