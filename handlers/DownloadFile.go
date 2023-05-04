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

	// Десериализация данных из запроса
	request_data := &protocol.FileName{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Поиск метаданных файла в базе данных
	file_metadata, err := db.GetFileMetadata(request.PublicKey, request_data.Name)
	if errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Файл с таким именем не существует"})
		return
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Открытие файла на чтение
	file, err := os.Open(fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_metadata.ID.Hex()))
	if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Расчёт количества блоков, на которые файл будет разбит в процессе передачи
	if stat, err := file.Stat(); err == nil {

		chunks := stat.Size() / config.PROTOCOL_CHUNK_SIZE

		if stat.Size()%config.PROTOCOL_CHUNK_SIZE != 0 {
			chunks += 1
		}

		protocol.SendResponse(conn, true, &protocol.FileMetadata{Name: file_metadata.Name, Encrypted: file_metadata.Encrypted, Chunks: uint32(chunks)})

	} else {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	}

	// Отправка всех блоков в сокет
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
