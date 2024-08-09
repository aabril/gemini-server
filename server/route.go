package server 

import (
	"net"
)

// HandlerFunc defines a function type for handling Gemini request
type HandlerFunc func(path stringg, conn net.Conn)

type Route struct {
	Path    string
	Handler HandlerFunc
}
