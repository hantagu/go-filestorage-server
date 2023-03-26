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

	// Find file's metadata in the database
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

		chunks := stat.Size() / config.PROTOCOL_CHUNK_SIZE

		if stat.Size()%config.PROTOCOL_CHUNK_SIZE != 0 {
			chunks += 1
		}

		protocol.SendResponse(conn, true, &protocol.FileMetadata{Name: file_metadata.Name, Encrypted: file_metadata.Encrypted, Chunks: uint32(chunks)})

	} else {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	}

	// Send all chunks of a file
	buffer := make([]byte, config.PROTOCOL_CHUNK_SIZE)
	for i := uint32(0); ; i++ {

		n, err := io.ReadFull(file, buffer)

		if n > 0 {
			protocol.SendResponse(conn, true, &protocol.FileChunk{Chunk: i, Content: buffer[:n]})
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

}
