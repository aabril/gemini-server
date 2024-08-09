package server

import (
	"fmt"
	"io"
	"net"
)

// WelcomeHandler handles the root path.
func WelcomeHandler(path string, conn net.Conn) {
	io.WriteString(conn, "20 text/gemini\r\nWelcome to the Gemini server!\r\n")
}

// HelloHandler handles the /hello path.
func HelloHandler(path string, conn net.Conn) {
	io.WriteString(conn, "20 text/gemini\r\nHello Gemini!\r\n")
}

// EchoHandler echoes back the requested path.
func EchoHandler(path string, conn net.Conn) {
	response := fmt.Sprintf("20 text/gemini\r\nYou Said: %s\r\n", path)
}
