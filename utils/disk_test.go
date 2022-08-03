package utils

import "testing"

func TestRetrieveSpareSpace(t *testing.T) {
	a, b, c := RetrieveSpareSpace()
	t.Logf("free MiB: %d, work mode: %d, writable: %v", a, b, c)
	if a < 0 || b < 0 || c == false {
		t.Fail()
	}
}
