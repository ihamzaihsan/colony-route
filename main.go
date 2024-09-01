package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
}
