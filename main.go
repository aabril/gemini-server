package main

import (
	"gemini-server/server"
	"log"
)

func main() {
	geminiServer := server.NewServer(":1965")

	// Add dynamic routes
	geminiServer.AddRoute("/", server.WelcomeHandler)
	geminiServer.AddRoute("/hello", server.HelloHandler)
	geminiServer.AddRoute("/echo", server.EchoHandler)

	if err := geminiServer.Start("server.crt", "server.key"); err != nil {
		log.Fatal(err)
	}
}
