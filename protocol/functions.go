package protocol

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"go-filestorage-server/config"
	"io"
	"net"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrPacketSignature = errors.New("неверная криптографическая подпись пакета")

func ReceiveAndVerifyPacket(conn net.Conn) (*Request, error) {

	buffer := make([]byte, config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE)

	// Чтение первых 4 байт (согласно документации BSON), которые указывают размер всего BSON документа
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return nil, err
	}

	// Конвертация байтов в значение типа UInt32
	packetLength := binary.LittleEndian.Uint32(buffer) - config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE

	// Чтение всего BSON документа
	buffer = append(buffer, make([]byte, packetLength)...)
	if _, err := io.ReadFull(conn, buffer[config.PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE:]); err != nil {
		return nil, err
	}

	// Десериализация BSON документа в структуру пакета общего типа
	packet := Request{}
	if err := bson.Unmarshal(buffer, &packet); err != nil {
		return nil, err
	}

	// Получение байтов подписи
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := io.ReadFull(conn, signature); err != nil {
		return nil, err
	}

	// Проверка криптографической подписи
	if !ed25519.Verify(packet.PublicKey, buffer, signature) {
		return nil, ErrPacketSignature
	}

	return &packet, nil
}

func SendResponse(conn net.Conn, successful bool, response_data interface{}) {

	raw_response_data, _ := bson.Marshal(response_data)

	raw_response, _ := bson.Marshal(&Response{
		Successful: successful,
		Data:       raw_response_data,
	})

	conn.Write(raw_response)
}
