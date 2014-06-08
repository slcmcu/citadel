package repository

import (
	"encoding/json"
	"path"
	"path/filepath"
	"strings"
)

// marshal encodes the value into a string via the json encoder
func (e *Repository) marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// unmarshal decodes the data using the json decoder into the value v
func (e *Repository) unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// isNotFoundErr returns true if the error is of type Key Not Found
func isNotFoundErr(err error) bool {
	return strings.Contains(err.Error(), "Key not found")
}

func buildServiceName(fullPath, prefix string) string {
	if fullPath == "" || fullPath == "/" {
		return "/"
	}

	dir, name := filepath.Split(fullPath)

	var (
		parts = strings.Split(dir, "/")
		full  = []string{}
	)

	switch len(parts) {
	case 0:
		return path.Join(name, prefix)
	}

	for _, p := range parts {
		if p != "" {
			full = append(full, p, "services")
		}
	}

	return path.Join(append(full, name, prefix)...)
}
