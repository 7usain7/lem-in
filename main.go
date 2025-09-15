package main

import (
	"fmt"
	"lem-in/funcs"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: go run main.go <File_Name>")
		return
	}

	if !strings.HasSuffix(args[0], ".txt") {
		fmt.Printf("Error: The file '%s' is not a .txt file.\n", args[0])
		return
	}

	var colony funcs.Colony
	colony.Setup()

	fileParsed, fileContent := funcs.ParseInput(args[0], &colony)
	if !fileParsed {
		return
	}

	// Display the input file content
	fmt.Print(fileContent)

}
