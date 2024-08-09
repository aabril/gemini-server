package tls

import (
	"crypto/tls"
	"fmt"
)

// LoadTLSConfig loads the TLS certficate and key.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	return config, nil
}
