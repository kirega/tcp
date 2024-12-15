package tcp

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	serverAddress    = "127.0.0.1:4000"
	clientRetryDelay = 5 * time.Second
	// TODO: Empirical Maximum number of clients on my machine, figure out how to get more connections and why the server was rejecting connections
	maxClients      = 150
	bufferSize      = 1024
	responseMessage = "Message received."
	clientMessage   = "Hello from client"
)

// Server represents the TCP server
type Server struct {
	listener         net.Listener
	totalConnections int
	connectionsLock  sync.Mutex
	newConnection    chan struct{}
}

// StartServer starts a TCP server on the specified port.
func StartServer() {
	server := &Server{
		newConnection: make(chan struct{}),
	}

	var err error
	server.listener, err = net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Error closing listener: %v", err)
		}
	}(server.listener)

	log.Printf("Server started on %s\n", serverAddress)

	go server.trackConnections()

	// Accept incoming connections
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go server.handleRequest(conn)
	}
}

// trackConnections increments and logs the number of active connections.
func (s *Server) trackConnections() {
	for range s.newConnection {
		s.connectionsLock.Lock()
		s.totalConnections++
		count := s.totalConnections
		s.connectionsLock.Unlock()
		log.Printf("New connection received. Total connections: %d\n", count)
	}
}

// handleRequest processes an incoming client connection.
func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	// Read the incoming data
	buf := make([]byte, bufferSize)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	log.Printf("Received: %s\n", string(buf[:n]))
	s.newConnection <- struct{}{}

	// Send a response to the client
	_, err = conn.Write([]byte(responseMessage))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// StartClient starts multiple TCP clients to send data to the server.
func StartClient() {
	var wg sync.WaitGroup

	for i := 0; i < maxClients; i++ {
		wg.Add(1)
		go sendData(&wg)
	}

	wg.Wait()
}

// sendData sends data to the TCP server and waits for a response.
func sendData(wg *sync.WaitGroup) {
	defer wg.Done()

	var conn net.Conn
	var err error

	// Retry until connection is established
	for conn == nil {
		conn, err = net.Dial("tcp", serverAddress)
		if err != nil {
			log.Printf("Error dialing server, retrying in %v: %v", clientRetryDelay, err)
			time.Sleep(clientRetryDelay)
		}
	}
	defer conn.Close()

	// Send a message to the server
	_, err = conn.Write([]byte(clientMessage))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	// Read response from the server
	buf := make([]byte, bufferSize)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	log.Printf("Received response: %s\n", string(buf[:n]))
}
