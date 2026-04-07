package main

import "fmt"

func VetSample() {
	// Printf argument mismatch for go vet to catch
	fmt.Printf("%d", "this should be an integer")
}
