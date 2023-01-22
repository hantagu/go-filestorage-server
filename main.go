package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var logger *log.Logger

func main() {

	// Initialize logger object
	logger = log.New(os.Stdout, "[FileSharing Server] ", log.Lmsgprefix|log.Ldate|log.Ltime)

	// Load server's TLS certificate
	tlsCert, err := tls.LoadX509KeyPair(TLS_CRT_PATH, TLS_KEY_PATH)
	catch(err)
	logger.Println("TLS certificate and key files read successfully")

	// Create a server's TLS configuration with loaded certificate and v1.2 <= TLS Version <= v1.3
	tlsConfig := tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}

	// Start a TLS listener
	listener, err := tls.Listen(LISTEN_NETWORK, LISTEN_ADDRESS, &tlsConfig)
	catch(err)
	logger.Printf("Started listening for incoming connections at %s\n", listener.Addr())

	// Create a channel to receive OS signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	// Create a WaitGroup that will wait for all connections currently being processed to complete
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			waitGroup.Wait()
			listener.Close()
		}
	}()

	// Endless loop that accepts new connections
	for {
		connection, err := listener.Accept()

		// This error is returned when the Listener is closed
		if errors.Is(err, net.ErrClosed) {
			logger.Println("Server successfully stopped")
			break
		} else if err != nil {
			logger.Printf("Accept() error: %s\n", err.Error())
			continue
		}

		logger.Printf("Accepted a new connection from %s\n", connection.RemoteAddr())

		// Add this connection to a WaitGroup
		waitGroup.Add(1)

		// Run a goroutine to hangle accepted connection
		handleConnection(connection, waitGroup)
	}
}

func catch(err error) {
	if err != nil {
		logger.Fatalln(err.Error())
	}
}
