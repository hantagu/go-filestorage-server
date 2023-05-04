package handlers

import (
	"bytes"
	"crypto/ed25519"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

func GrantAccess(conn net.Conn, request *protocol.Request) {

	// Десериализация данных из запроса
	request_data := &protocol.FileAccess{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	if len(request_data.PublicKey) != ed25519.PublicKeySize {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Неверная длина публичного ключа"})
		return
	} else if bytes.Equal(request.PublicKey, request_data.PublicKey) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Вы не можете предоставить доступ к своему файлу самому себе"})
		return
	}

	// Вставка публичного ключа в метаданные файла в базе данных
	if matched, modified, err := db.GrantAccess(request.PublicKey, request_data.Name, request_data.PublicKey); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	} else if matched == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Файл с таким именем не существует"})
	} else if modified == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Этот пользователь уже имеет доступ к этому файлу"})
	} else {
		protocol.SendResponse(conn, true, &protocol.Empty{})
	}
}
