package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

// GeminiServer holds the TLS listener and routes.
type GeminiServer struct {
	Addr   string
	Routes []Route
}

// NewServer initializes a new Gemini server with a specified address.
func NewServer(addr string) *GeminiServer {
	return &GeminiServer{
		Addr:   addr,
		Routes: []Route{},
	}
}

// AddRoute adds a new route to the Gemini Server.
func (s *GeminiServer) AddRoute(path string, handler HandlerFunc) {
	s.Routes = append(s.Routes, Route{Path: path, Handler: handler})
}

// HandleRequest processes incoming Gemini requests and involes the corresponding route handler
func (s *GeminiServer) HandleRequest(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Failed to read from connection:", err)
		return
	}

	// Trim whitespace and newline characters
	path := strings.TrimSpace(string(buf[:n]))
	log.Println("Received request for:", path)

	// Normalise the path to remove trailing slashes
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	for _, route := range s.Routes {
		if strings.HasPrefix(path, route.Path) {
			route.Handler(path, conn)
			return
		}
	}

	// If no route matches, send a 51 "not found" response.
	conn.Write([]byte("51 Not Found\r\n"))
}

func (s *GeminiServer) Start(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", s.Addr, config)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	log.Printf("Gemini server is listening on %s", s.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		go s.HandleRequest(conn)
	}
}
