package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
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

// Return the value in /etc/machine-id
func GetMachineID() (string, error) {
	data, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("host does not support /etc/machine-id")
		}
		return "", err
	}
	return strings.Trim(string(data), "\n"), nil
}

func RemoveTag(name string) string {
	if hasTag, index := checkTag(name); hasTag {
		return name[:index]
	}
	return name
}

func checkTag(name string) (bool, int) {
	index := strings.LastIndex(name, ":")
	if index == -1 || strings.Contains(name[index:], "/") {
		return false, -1
	}
	return true, index
}
