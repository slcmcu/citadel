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

// Return the value in /etc/hostname
func GetMachineID() (string, error) {
	data, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("host does not support /etc/machine-id")
		}
		return "", err
	}
	return strings.Trim(string(data), "\n"), nil
}

func checkTag(name string) (bool, int) {
	index := strings.LastIndex(name, ":")
	if index == -1 || strings.Contains(name[index:], "/") {
		return false, -1
	}
	return true, index
}

func RemoveTag(name string) string {
	if hasTag, index := checkTag(name); hasTag {
		return name[:index]
	}
	return name
}

func RemoveSlash(name string) string {
	return strings.Replace(name, "/", "", -1)
}

func SplitURI(uri string) (string, string) {
	arr := strings.Split(uri, "://")
	if len(arr) == 1 {
		return "unix", arr[0]
	}
	prot := arr[0]
	if prot == "http" {
		prot = "tcp"
	}
	return prot, arr[1]
}

func CleanImageName(name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 1 {
		return RemoveSlash(RemoveTag(name))
	}
	return CleanImageName(parts[1])
}
