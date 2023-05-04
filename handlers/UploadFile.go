package handlers

import (
	"fmt"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/protocol"
	"net"
	"os"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadFile(conn net.Conn, request *protocol.Request) {

	// Генерация ID файла и пути, по которому он будет сохранён
	file_id := primitive.NewObjectID()
	file_path := fmt.Sprintf("%s%c%s", config.Config.UserdataPath, os.PathSeparator, file_id.Hex())

	// Создание файла
	file, err := os.Create(file_path)
	if err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		return
	}
	defer file.Close()

	// Десериализация данных из запроса
	request_data := &protocol.FileMetadata{}
	if err := bson.Unmarshal(request.Data, request_data); err != nil {
		protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		goto exit
	} else if request_data.Chunks == 0 {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Файл не может быть пустым"})
		goto exit
	} else if !regexp.MustCompile(`^[^[:cntrl:]/]+$`).MatchString(request_data.Name) {
		protocol.SendResponse(conn, false, &protocol.Description{Description: "Недопустимое имя файла"})
		goto exit
	}

	// Попытка вставить метаданные файла в базу данных
	if err := db.InsertFileMetadata(file_id, request.PublicKey, request_data.Name, request_data.Encrypted); err != nil {

		if mongo.IsDuplicateKeyError(err) {
			protocol.SendResponse(conn, false, &protocol.Description{Description: "Файл с таким именем уже существует"})
		} else {
			protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
		}

		goto exit
	}

	// Если метаданные были успешно сохранены в базе данных, тогда отправляем ответ о готовности принять сам файл
	protocol.SendResponse(conn, true, &protocol.Empty{})

	// Получение всех блоков файла
	for i := uint32(0); i < request_data.Chunks; i++ {

		request, err := protocol.ReceiveAndVerifyPacket(conn)

		if err != nil || request.Type != protocol.REQ_UPLOAD_FILE {

			if err != nil {
				protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
			} else {
				protocol.SendResponse(conn, false, &protocol.Description{Description: "Неверный запрос"})
			}

			goto exit
		}

		request_data := &protocol.FileChunk{}
		if err := bson.Unmarshal(request.Data, request_data); err != nil {
			protocol.SendResponse(conn, false, &protocol.Description{Description: err.Error()})
			goto exit
		}

		file.Write(request_data.Content)
		protocol.SendResponse(conn, true, &protocol.Empty{})
	}

	return

exit:
	os.Remove(file_path)
	db.DeleteFileMetadataByID(file_id)
}
