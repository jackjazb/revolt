package main

import (
	"fmt"
	"os"
)

func run() error {
	RunServer()
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
