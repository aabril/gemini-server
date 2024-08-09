package server

import (
	"net"
)

// HandlerFunc defines a function type for handling Gemini request
type HandlerFunc func(path string, conn net.Conn)

// Define the RouteMatcher type
type RouteMatcher func(path string, route Route) bool

// Define the Route struct
type Route struct {
	Path    string
	Handler HandlerFunc
}
