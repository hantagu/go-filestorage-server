package main

import (
	"crypto/tls"
	"errors"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

func main() {

	log.Default().SetFlags(0)

	config.Init()

	if err := db.Init(); err != nil {
		log.Default().Fatal("Не удалось подключиться к MongoDB\n\nПопробуйте запустить MongoDB в существующем Docker-контейнере:\n\tdocker start filestorage_db\n\nИли запустите новый:\n\tdocker run -d -p 127.0.0.1:27017:27017 --name filestorage_db mongo\n\n")
	}

	// Создание директории, в которой будут храниться все файлы пользователей, если она не существует
	if err := os.Mkdir(config.Config.UserdataPath, 0o777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Default().Fatalln(err)
	}

	// Загрузка TLS сертификата сервера
	tlsCert, err := tls.LoadX509KeyPair(config.Config.TLSCertificatePath, config.Config.TLSKeyPath)
	if err != nil {
		log.Default().Fatalln(err)
	}

	// Создание конфигурации TLS сервера с загруженным сертификатом и минимальной версией протокола TLS 1.2
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}

	// Запуск TLS сервера
	listener, err := tls.Listen("tcp", config.Config.ListenAddress, tlsConfig)
	if err != nil {
		log.Default().Fatalln(err)
	}

	log.Default().Printf("Сервер был запущен и ожидает подключений по адресу %s\nНажмите Ctrl-C для завершения работы", listener.Addr())

	// Создание канала для получения сигналов из операционной системы
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	// Создание объекта WaitGroup, который будет ожидать завершения всех соединений
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			log.Default().Println("\nОжидание завершения всех соединений...\nНажмите Ctrl-C снова, чтобы завершить работу принудительно")
			go func() {
				for range shutdownChan {
					log.Default().Println("\nПринудительное завершение работы")
					listener.Close()
					os.Exit(0)
				}
			}()
			waitGroup.Wait()
			listener.Close()
			os.Exit(0)
		}
	}()

	// Бесконечный цикл, принимающий новые соединения
	for {
		conn, err := listener.Accept()

		// Эта ошибка возвращается, когда сервер прекращает принимать новые соединения
		if errors.Is(err, net.ErrClosed) {
			log.Default().Println("Сервер успешно остановлен")
			break
		} else if err != nil {
			log.Default().Printf("%s: Ошибка при принятии нового соединения: %s\n", conn.RemoteAddr(), err)
			continue
		}

		log.Default().Printf("Принято новое соединение от %s\n", conn.RemoteAddr())

		// Добавление этого соединения к "счётчику"
		waitGroup.Add(1)

		// Запуск горутины, которая будет обрабатывать это соединение
		go handleConnection(conn, waitGroup)
	}
}
