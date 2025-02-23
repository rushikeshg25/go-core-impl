package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	listener   net.Listener
	address    string
	messages   chan []byte
	clients    map[net.Conn]bool
	clientsMux sync.RWMutex
	done       chan struct{}
}

func NewServer(address string, bufferSize int) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener:   listener,
		address:    address,
		messages:   make(chan []byte, bufferSize),
		clients:    make(map[net.Conn]bool),
		clientsMux: sync.RWMutex{},
		done:       make(chan struct{}),
	}, nil
}

func (s *Server) Start() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		s.Stop()
	}()

	defer s.listener.Close()
	go s.broadcast()

	log.Printf("Server started on %s", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil
			default:
				log.Printf("Error accepting connection: %v", err)
				continue
			}
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {
	log.Println("Shutting down server...")
	close(s.done)
	s.listener.Close()

	s.clientsMux.Lock()
	for conn := range s.clients {
		conn.Close()
	}
	s.clientsMux.Unlock()

	close(s.messages)
	log.Println("Server shutdown complete")
}

func (s *Server) broadcast() {
	for msg := range s.messages {
		s.clientsMux.RLock()
		for client := range s.clients {
			writer := bufio.NewWriter(client)
			_, err := writer.Write(msg)
			if err == nil {
				_, err = writer.Write([]byte{'\n'})
			}
			if err == nil {
				err = writer.Flush()
			}
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}
		s.clientsMux.RUnlock()
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New client connected: %s", conn.RemoteAddr())

	s.clientsMux.Lock()
	s.clients[conn] = true
	s.clientsMux.Unlock()

	defer func() {
		s.clientsMux.Lock()
		delete(s.clients, conn)
		s.clientsMux.Unlock()
		log.Printf("Client disconnected: %s", conn.RemoteAddr())
	}()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Error reading from client %s: %v", conn.RemoteAddr(), err)
			return
		}

		message = message[:len(message)-1]

		if len(message) > 0 {
			select {
			case s.messages <- message:
				log.Printf("Received message from %s: %s", conn.RemoteAddr(), string(message))
			case <-s.done:
				return
			}
		}
	}
}

func main() {
	server, err := NewServer(":8080", 1024)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	log.Println("Starting server...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
