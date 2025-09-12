package main

import (
	"fmt"
	"lem-in/funcs"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	colony, err := funcs.ParseInput(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Print test.txt file content
	funcs.DisplayColony(colony)
}
