package utils

import "testing"

func TestToCpuSet(t *testing.T) {
	expected := "0,1"

	actual := IToCpuset(2)

	if actual != expected {
		t.Fatalf("expected %s received %s", expected, actual)
	}
}

func TestCpusetToI(t *testing.T) {

	actual := CpusetTOI("0,2,1")
	if actual != 3 {
		t.Fatalf("expected 3 cores received %d", actual)
	}
}
