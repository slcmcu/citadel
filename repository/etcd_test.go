package repository

import "testing"

func TestBuildNameWithSlash(t *testing.T) {
	expected := "/"

	name := buildServiceName("/", "services")

	if name != expected {
		t.Fatalf("%s != %s", name, expected)
	}
}

func TestBuildNameWithSingleName(t *testing.T) {
	expected := "local/services"

	name := buildServiceName("local", "services")

	if name != expected {
		t.Fatalf("%s != %s", name, expected)
	}
}
