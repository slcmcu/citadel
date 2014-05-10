package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateUUID(size int) string {
	id := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic(err)
	}
	return hex.EncodeToString(id)
}
