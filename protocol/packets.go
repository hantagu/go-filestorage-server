package protocol

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
)

// Общий пакет запросов
type Request struct {
	Type      string            `bson:"type"`       // Тип запроса
	PublicKey ed25519.PublicKey `bson:"public_key"` // Публичный ключ отправителя
	Data      bson.Raw          `bson:"data"`       // Данные запроса
}

// Общий пакет ответов
type Response struct {
	Successful bool     `bson:"successful"` // Был ли запрос успешным?
	Data       bson.Raw `bson:"data"`       // Данные ответа
}

// Пустой пакет
type Empty struct{}

// Данные, содержащие строковое описание. Используется в ответах от сервера
type Description struct {
	Description string `bson:"description"`
}

// -------------------------------------------------- Пользователи -------------------------------------------------- //

// Данные, содержащие публичный ключ
type PublicKey struct {
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Данные, содержащие имя пользователя
type Username struct {
	Username string `bson:"username"`
}

// -------------------------------------------------- Файлы -------------------------------------------------- //

// Данные, содержащие название файла
type FileName struct {
	Name string `bson:"name"`
}

// Данные, содержащие название файла и публичный ключ, которому будет предоставлен / отозван доступ к указанному файлу
type FileAccess struct {
	Name      string            `bson:"name"`
	PublicKey ed25519.PublicKey `bson:"public_key"`
}

// Данные, содержащие один блок файла, на которые он был разбит
type FileChunk struct {
	Chunk   uint32 `bson:"chunk"`   // Номер этого блока
	Content []byte `bson:"content"` // Его содержимое
}

// Данные, содержащие метаданные файла
type FileMetadata struct {
	Name      string `bson:"name"`      // Название файла
	Encrypted bool   `bson:"encrypted"` // Был ли файл дополнительно зашифрован на стороне клиента?
	Chunks    uint32 `bson:"chunks"`    // Общее количество блоков
}
