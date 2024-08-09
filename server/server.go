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

// HandleRequest helpers
func exactMatch(path string, route Route) bool {
	return path == route.Path
}

func prefixMatch(path string, route Route) bool {
	hasPrefix := strings.HasPrefix(path, route.Path)
	lenPathMatches := len(path) == len(route.Path)
	lenRoutePathMatches := len(path) > len(route.Path) && path[len(route.Path)] == '/'
	fmt.Printf("Checking prefix match: %s starts with %s? %v\n", path, route.Path, hasPrefix) //TODO remove
	fmt.Printf("Length match: %v, Slash match: %v\n", lenPathMatches, lenRoutePathMatches)    // TODO remove
	return hasPrefix && (lenPathMatches || lenRoutePathMatches)
}

func matchRoute(routes []Route, matcher RouteMatcher) func(string) *Route {
	return func(path string) *Route {
		for _, route := range routes {
			if matcher(path, route) {
				return &route
			}
		}
		return nil
	}
}

func handleNotFound(conn net.Conn) {
	conn.Write([]byte("51 Not Found \r\n"))
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

	path := strings.TrimSpace(string(buf[:n]))
	log.Println("Received request for:", path)

	// Extract the path from the URI (e.g., "gemini://localhost/somepath" -> "/somepath")
	if strings.HasPrefix(path, "gemini://") {
		path = path[len("gemini://"):]
		// Remove host part
		if idx := strings.Index(path, "/"); idx != -1 {
			path = path[idx:]
		} else {
			path = "/"
		}
	}

	// Normalize the path to remove query strings and trailing slashes
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	// Normalize the path to remove trailing slashes except for the root "/"
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Create a matcher function that prioritizes exact match and falls back to prefix match
	matcher := func(path string) *Route {
		if route := matchRoute(s.Routes, exactMatch)(path); route != nil {
			return route
		}
		return matchRoute(s.Routes, prefixMatch)(path)
	}

	// Find the matching route and execute the handler
	if route := matcher(path); route != nil {
		// log.Println("Route matched:", route.Path)
		route.Handler(path, conn)
	} else {
		// log.Println("No matching route found.")
		handleNotFound(conn)
	}
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
