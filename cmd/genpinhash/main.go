// cmd/genpinhash/main.go
// Run: go run ./cmd/genpinhash <your-pin>
// Copy the output into your .env as PIN_HASH=<value>
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/genpinhash <pin>")
		os.Exit(1)
	}
	pin := os.Args[1]
	h := sha256.Sum256([]byte(pin))
	fmt.Printf("PIN_HASH=%s\n", hex.EncodeToString(h[:]))
}
