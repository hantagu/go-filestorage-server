package config

const (
	CONFIG_FILE_PATH                   = "./config.json"                    // Путь до файла конфигурации по умолчанию
	PROTOCOL_PREAMBLE                  = "\xB3\xD6\x4C\xA4\xF6\x9B\x71\xF8" //
	PROTOCOL_PREAMBLE_SIZE             = 8                                  //
	PROTOCOL_CHUNK_SIZE                = 8 * 1024 * 1024                    // Размер одного блока файла (8 МиБ)
	PROTOCOL_BSON_DOCUMENT_LENGTH_SIZE = 4                                  // Количество байт, в которых указан размер всего BSON документа
	MONGODB_CONTEXT_TIMEOUT            = 3                                  // Время ожидания (в секундах) ответов от MongoDB (передаётся в context'ы)
)
