package main

import "testing"

func TestMain(t *testing.T) {
	r := []rune("4325")
	if n, err := strtol(&r); err != nil {
		t.Error(err)
	} else if n != 4325 {
		t.Fatalf("expected %d, but got %d\n", 4325, n)
	} else if len(r) != 0 {
		t.Fatalf("expected %d, but got %d\n", 0, len(r))
	}
}
