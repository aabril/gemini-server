# Gemini Server

An implementation of the gemini protocol as a server using golang.

## Code 

- `GeminiServer struct` contains the address and routes for the server.
- `NewServer` initializes a new server
- `AddRoute` allows adding routes dynamically with corresponding handlers
- `HandleRequest` processes incoming request and dispatches them to the appropriate ahndler based on the path.
- `Start` starts the TLS server, listens for connections and handles them concurrently

## Create TLS certficates

`openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt`


