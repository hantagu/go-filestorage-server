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
	"syscall"
)

func main() {

	log.Default().SetFlags(0)

	config.Init()

	if err := db.Init(); err != nil {
		log.Default().Fatal("Failed to connect to MongoDB\n\nTry running an existing Docker container with MongoDB:\n\tdocker start -ai filestorage_db\n\nOr create a new Docker container with MongoDB:\n\tdocker run -p 127.0.0.1:27017:27017 --name filestorage_db mongo\n\n")
	}

	// Create User Data directory if it doesn't exist
	if err := os.Mkdir(config.Config.UserdataPath, 0o777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Default().Fatalln(err)
	}

	// Load a TLS server certificate
	tlsCert, err := tls.LoadX509KeyPair(config.Config.TLSCertificatePath, config.Config.TLSKeyPath)
	if err != nil {
		log.Default().Fatalln(err)
	}

	// Create a TLS server configuration with the loaded certificate and the minimum TLS 1.2 version
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}

	// Start a TLS listener
	listener, err := tls.Listen("tcp", config.Config.ListenAddress, tlsConfig)
	if err != nil {
		log.Default().Fatalln(err)
	}

	log.Default().Printf("The listener has been started on %s\n", listener.Addr())

	// Create a channel to receive signals from OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	// Create a WaitGroup that will wait for all connections currently being processed to complete
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			log.Default().Println("\nWaiting for the completion of the current connections.\nPress Ctrl-C again to force shutdown")
			go func() {
				for range shutdownChan {
					log.Default().Println("\nForced shutdown")
					listener.Close()
					os.Exit(0)
				}
			}()
			waitGroup.Wait()
			listener.Close()
			os.Exit(0)
		}
	}()

	// Endless loop that accepts new connections
	for {
		conn, err := listener.Accept()

		// This error is returned when the Listener is closed
		if errors.Is(err, net.ErrClosed) {
			log.Default().Println("Server successfully stopped")
			break
		} else if err != nil {
			log.Default().Printf("Connection accept error: %s\n", err)
			continue
		}

		log.Default().Printf("Accepted a new connection from %s\n", conn.RemoteAddr())

		// Add this connection to a WaitGroup
		waitGroup.Add(1)

		// Run a goroutine to handle accepted connection
		go handleConnection(conn, waitGroup)
	}
}
