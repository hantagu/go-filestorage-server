package db_types

import (
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID  `bson:"_id"`       // Уникальный ID файла, генерируемый MongoDB
	Owner     ed25519.PublicKey   `bson:"owner"`     // Публичный ключ владельца файла
	Name      string              `bson:"name"`      // Имя файла
	Encrypted bool                `bson:"encrypted"` // Был ли файл дополнительно зашифрован на стороне клиента?
	Access    []ed25519.PublicKey `bson:"access"`    // Список публичных ключей других пользователей, которым разрешён доступ на чтение к этому файлу
}
