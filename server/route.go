package server 

import (
	"net"
)

// HandlerFunc defines a function type for handling Gemini request
type HandlerFunc func(path string, conn net.Conn)

type Route struct {
	Path    string
	Handler HandlerFunc
}
