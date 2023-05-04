package handlers

import (
	"errors"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(conn net.Conn, request *protocol.Request) {

	request_data_publickey := &protocol.PublicKey{}
	request_data_username := &protocol.Username{}

	// Попытка десериализовать данные из запроса как публичный ключ
	err1 := bson.Unmarshal(request.Data, request_data_publickey)
	if err1 != nil {

		// Если не получилось, попытка десериализовать данные из запроса как имя пользователя
		err2 := bson.Unmarshal(request.Data, request_data_username)
		if err2 != nil {
			protocol.SendResponse(conn, false, &protocol.Description{Description: "Неверные данные в запросе"})
			return
		}

		if result, err := db.GetUserByUsername(request_data_username.Username); errors.Is(err, mongo.ErrNoDocuments) {
			protocol.SendResponse(conn, false, &protocol.Description{Description: "Пользователь с таким именем не найден"})
			return
		} else if err != nil {
			protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
			return
		} else {
			protocol.SendResponse(conn, true, &protocol.PublicKey{PublicKey: result})
			return
		}
	}

	if result, err := db.GetUserByPublicKey(request.PublicKey); errors.Is(err, mongo.ErrNoDocuments) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Пользователь с таким публичным ключом не найден"})
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	} else {
		protocol.SendResponse(conn, true, &protocol.Username{Username: result})
	}
}
