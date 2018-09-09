package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	var want bool
	var got bool

	want = true
	got = true

	if got != want {
		t.Fatalf("Got '%v', expected '%v'", got, want)
	}
}
