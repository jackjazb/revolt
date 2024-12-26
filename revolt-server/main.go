package main

import (
	"fmt"
	"os"
)

func run() error {
	err := RunServer()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
