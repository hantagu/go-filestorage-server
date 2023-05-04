package handlers

import (
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUsername(conn net.Conn, request *protocol.Request) {

	// Десериализация данных из запроса
	request_data := &protocol.Username{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}

	// Проверка допустимости имени пользователя при помощи регулярного выражения
	if !regexp.MustCompile(`^[a-z0-9_]{5,}$`).MatchString(request_data.Username) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Вы можете использовать только символы a-z, 0-9 и символ подчеркивания, также минимальная длина имени пользователя должна быть не менее 5 символов"})
		return
	}

	// Попытка изменить имя пользователя и сохранить изменения в базе данных
	if err := db.SetUsername(request.PublicKey, request_data.Username); mongo.IsDuplicateKeyError(err) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Это имя пользователя уже занято"})
	} else if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
	} else {
		protocol.SendResponse(conn, true, &protocol.Empty{})
	}
}
