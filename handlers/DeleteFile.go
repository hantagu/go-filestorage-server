package handlers

import (
	"errors"
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteFile(conn net.Conn, request *protocol.Request) {

	// Десериализация данных из запроса
	request_data := &protocol.FileName{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Удаление метаданных файла из базы данных
	file_metadata, err := db.DeleteFileMetadata(request.PublicKey, request_data.Name)

	if errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Файл с таким именем не существует"})
		return
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Удаление самого файла
	os.Remove(fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_metadata.ID.Hex()))

	// Отправка ответа
	protocol.SendResponse(conn, true, &protocol.Empty{})
}
