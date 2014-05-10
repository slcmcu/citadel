package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strings"
)

var (
	ErrUnableToGenerateUUID = errors.New("unable to generate uuid from interfaces")
)

// GenerateUUID generates a random id with the specified size
func GenerateUUID(size int) string {
	id := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic(err)
	}
	return hex.EncodeToString(id)
}

// GetUUID uses the mac address for the external interface to generate
// a uuid
func GetUUID() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if n := iface.HardwareAddr.String(); n != "" {
			return strings.Replace(n, ":", "", -1), nil
		}
	}
	return "", ErrUnableToGenerateUUID
}
