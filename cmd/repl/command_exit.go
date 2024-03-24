package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *config, args ...string) error {
	fmt.Println("goodbye!")
	os.Exit(0)
	return nil
}
