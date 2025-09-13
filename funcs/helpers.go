package funcs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseInput(inputFile string, colony *Colony) (bool, string) {
	hasProcessedFirstLine := false
	inputData, openErr := os.Open(inputFile)
	if openErr != nil {
		fmt.Println("ERROR: unable to open file", openErr)
		os.Exit(1)
	}
	defer inputData.Close()

	fileContent := ""
	currentLineNumber := 0
	scanner := bufio.NewScanner(inputData)
	isStart := false
	isEnd := false
	totalAnts := 0
	startRoomDefined := false
	endRoomDefined := false

	for scanner.Scan() {
		currentLine := scanner.Text()
		fileContent += currentLine + "\n"
		currentLineNumber++

		// Process comments and special directives
		if strings.HasPrefix(currentLine, "#") {
			if currentLine == "##start" {
				if startRoomDefined {
					fmt.Printf("ERROR: invalid data format, multiple start rooms defined (%s)\n", currentLine)
					os.Exit(1)
				}
				startRoomDefined = true
				isStart = true
				continue
			} else if currentLine == "##end" {
				if endRoomDefined {
					fmt.Printf("ERROR: invalid data format, multiple end rooms defined (%s)\n", currentLine)
					os.Exit(1)
				}
				endRoomDefined = true
				isEnd = true
				continue
			} else {
				continue
			}
		}

		// Parse initial line (ant count)
		if !hasProcessedFirstLine {
			lineComponents := strings.Fields(currentLine)
			if len(lineComponents) != 1 {
				fmt.Printf("Invalid data format detected (%s)\n", currentLine)
				os.Exit(1)
			}
			totalAnts, openErr = strconv.Atoi(lineComponents[0])
			if openErr != nil || totalAnts <= 0 {
				fmt.Printf("ERROR: invalid data format, invalid number of ants (%s)\n", currentLine)
				os.Exit(1)
			}
			hasProcessedFirstLine = true
			continue
		}
	}
}
