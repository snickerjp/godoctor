package main

import "testing"

func TestHello(t *testing.T) {
	// Simple test
	if 1+1 != 2 {
		t.Error("1+1 should be 2")
	}
}
