package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strconv"
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

func IToCpuset(cpus []int) string {
	return strings.Join(toA(cpus), ",")
}

func toA(i []int) []string {
	s := make([]string, len(i))

	for j, ii := range i {
		s[j] = strconv.Itoa(ii)
	}

	return s
}

func CpusetTOI(cpuset string) []int {
	var (
		s = strings.Split(cpuset, ",")
		i = make([]int, len(s))
	)

	for j, ss := range s {
		si, err := strconv.Atoi(ss)
		if err != nil {
			panic(err)
		}

		i[j] = si
	}

	return i
}
