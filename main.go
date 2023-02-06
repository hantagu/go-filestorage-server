package main

import (
	"crypto/tls"
	"errors"
	"go-filestorage-server/mongodb"
	"go-filestorage-server/utils"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	utils.InitConfig()
	mongodb.InitMongoDB()

	// Create User Data directory if it doesn't exist
	if stat, err := os.Stat(utils.Config.UserdataPath); err != nil {
		if err := os.Mkdir(utils.Config.UserdataPath, 0o700); err != nil {
			utils.Logger.Fatalln(err)
		}
	} else if !stat.IsDir() {
		utils.Logger.Fatalf("`%s` is not a directory!\n", utils.Config.UserdataPath)
	}

	// Load a TLS server certificate
	tlsCert, err := tls.LoadX509KeyPair(utils.Config.TLSCertificatePath, utils.Config.TLSKeyPath)
	if err != nil {
		utils.Logger.Fatalln(err)
	}

	// Create a TLS server configuration with the loaded certificate and the minimum TLS 1.2 version
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{tlsCert},
	}

	// Start a TLS listener
	listener, err := tls.Listen("tcp", utils.Config.ListenAddress, tlsConfig)
	if err != nil {
		utils.Logger.Fatalln(err)
	}

	utils.Logger.Printf("The listener has been started on %s\n", listener.Addr())

	// Create a channel to receive signals from OS
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	// Create a WaitGroup that will wait for all connections currently being processed to complete
	waitGroup := &sync.WaitGroup{}

	go func() {
		for range shutdownChan {
			utils.Logger.Println("\nWaiting for the completion of the current connections.\nPress Ctrl-C again to force shutdown")
			go func() {
				for range shutdownChan {
					utils.Logger.Println("\nForced shutdown")
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
			utils.Logger.Println("Server successfully stopped")
			break
		} else if err != nil {
			utils.Logger.Printf("Connection accept error: %s\n", err)
			continue
		}

		utils.Logger.Printf("Accepted a new connection from %s\n", conn.RemoteAddr())

		// Add this connection to a WaitGroup
		waitGroup.Add(1)

		// Run a goroutine to handle accepted connection
		go handleConnection(conn, waitGroup)
	}
}
