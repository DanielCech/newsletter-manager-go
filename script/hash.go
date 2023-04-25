package main

import (
	"errors"
	"fmt"
	"os"

	"strv-template-backend-go-api/crypto"
)

const expectedArgs = 2

func die(err error) {
	_, err = fmt.Printf("error: %v\n", err)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) != expectedArgs {
		die(errors.New("usage: hash.go <password>"))
	}

	passwd := os.Args[1]
	pepper, ok := os.LookupEnv("HASH_PEPPER")
	if !ok {
		if _, err := fmt.Println("warning: env var HASH_PEPPER is empty"); err != nil {
			die(err)
		}
	}
	hasher := crypto.NewDefaultBcryptHasher([]byte(pepper))

	h, err := hasher.HashPassword([]byte(passwd))
	if err != nil {
		die(fmt.Errorf("hashing password: %w", err))
	}

	if _, err = fmt.Print(string(h)); err != nil {
		panic(err)
	}
}
