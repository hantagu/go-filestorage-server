package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initialize logger object
var logger *log.Logger = log.New(os.Stdout, "", 0)

func main() {

	// Create User Data directory if it doesn't exist
	if stat, err := os.Stat(USERDATA_PATH); err != nil || !stat.IsDir() {
		if err := os.Mkdir(USERDATA_PATH, USERDATA_DEFAULT_PERM); err != nil {
			catch(err)
		}
	}

	// Connect to MongoDB
	dbClient, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	catch(err)
	if dbClient.Connect(context.TODO()) != nil {
		catch(err)
	}

	// Load a TLS server certificate
	tlsCert, err := tls.LoadX509KeyPair(TLS_CRT_PATH, TLS_KEY_PATH)
	catch(err)

	// Create a TLS server configuration with the loaded certificate and the minimum TLS 1.2 version
	tlsConfig := tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}

	// Start a TLS listener
	listener, err := tls.Listen(LISTEN_NETWORK, LISTEN_ADDRESS, &tlsConfig)
	catch(err)
	logger.Printf("The listener has been started to handle incoming connections on %s\n", listener.Addr())

	// Create a channel to receive signals from OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	// Create a WaitGroup that will wait for all connections currently being processed to complete
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			logger.Println("\nWaiting for the completion of the current connections.\nPress Ctrl-C again to force shutdown")
			go func() {
				for range shutdownChan {
					logger.Println("\nForced shutdown")
					listener.Close()
					os.Exit(0)
				}
			}()
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

		// Run a goroutine to handle accepted connection
		handleConnection(connection, waitGroup)
	}
}

func catch(err error) {
	if err != nil {
		logger.Fatalln(err.Error())
	}
}
