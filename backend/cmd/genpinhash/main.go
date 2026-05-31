// cmd/genpinhash/main.go
// Run: go run ./cmd/genpinhash <your-pin>
// Copy the output into your .env as PIN_HASH=<value>
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(
			os.Stderr,
			"Usage: go run ./cmd/genpinhash <pin>",
		)
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(os.Args[1]),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"PIN_HASH=%s\n",
		string(hash),
	)
}
