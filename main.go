package main

import (
	"crypto/tls"
	"errors"
	"go-filestorage-server/config"
	"go-filestorage-server/db"
	"go-filestorage-server/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	config.Init()
	db.Init()

	// Create User Data directory if it doesn't exist
	if err := os.Mkdir(config.Config.UserdataPath, 0o700); !errors.Is(err, os.ErrExist) {
		logger.Logger.Fatalln(err)
	}

	// Load a TLS server certificate
	tlsCert, err := tls.LoadX509KeyPair(config.Config.TLSCertificatePath, config.Config.TLSKeyPath)
	if err != nil {
		logger.Logger.Fatalln(err)
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
		logger.Logger.Fatalln(err)
	}

	logger.Logger.Printf("The listener has been started on %s\n", listener.Addr())

	// Create a channel to receive signals from OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	// Create a WaitGroup that will wait for all connections currently being processed to complete
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			logger.Logger.Println("\nWaiting for the completion of the current connections.\nPress Ctrl-C again to force shutdown")
			go func() {
				for range shutdownChan {
					logger.Logger.Println("\nForced shutdown")
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
			logger.Logger.Println("Server successfully stopped")
			break
		} else if err != nil {
			logger.Logger.Printf("Connection accept error: %s\n", err)
			continue
		}

		logger.Logger.Printf("Accepted a new connection from %s\n", conn.RemoteAddr())

		// Add this connection to a WaitGroup
		waitGroup.Add(1)

		// Run a goroutine to handle accepted connection
		go handleConnection(conn, waitGroup)
	}
}
