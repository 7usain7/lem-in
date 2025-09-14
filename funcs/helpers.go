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
		fmt.Printf("ERROR: unable to open file '%s'. Details: %v\n", inputFile, openErr)
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
					fmt.Printf("ERROR: invalid data format, multiple ##start rooms defined (%s)\n", currentLine)
					os.Exit(1)
				}
				startRoomDefined = true
				isStart = true
				continue
			} else if currentLine == "##end" {
				if endRoomDefined {
					fmt.Printf("ERROR: invalid data format, multiple ##end rooms defined (%s)\n", currentLine)
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
				fmt.Printf("ERROR: invalid data format. Expected a single number for ant count, but found '%s'.\n", currentLine)
				os.Exit(1)
			}
			totalAnts, openErr = strconv.Atoi(lineComponents[0])
			if openErr != nil || totalAnts <= 0 {
				fmt.Printf("ERROR: invalid data format, invalid number of ants (%s). Must be a positive integer.\n", currentLine)
				os.Exit(1)
			}
			hasProcessedFirstLine = true
			continue
		}

		// Process start room specification
		if isStart {
			isStart = false
			roomData := strings.Fields(currentLine)
			if len(roomData) != 3 {
				fmt.Printf("ERROR: invalid data format, invalid room definition (%s)\n", currentLine)
				os.Exit(1)
			}
			roomName := roomData[0]
			xPos, xErr := strconv.Atoi(roomData[1])
			yPos, yErr := strconv.Atoi(roomData[2])
			if xErr != nil || yErr != nil {
				fmt.Printf("ERROR: invalid data format, invalid room definition (%s)\n", currentLine)
				os.Exit(1)
			}
			colony.CreateRoom(roomName, "start", xPos, yPos)
			continue
		}

		// Process end room specification
		if isEnd {
			isEnd = false
			roomData := strings.Fields(currentLine)
			if len(roomData) != 3 {
				fmt.Printf("ERROR: invalid data format, invalid room definition (%s)\n", currentLine)
				os.Exit(1)
			}
			roomName := roomData[0]
			xPos, xErr := strconv.Atoi(roomData[1])
			yPos, yErr := strconv.Atoi(roomData[2])
			if xErr != nil || yErr != nil {
				fmt.Printf("ERROR: invalid data format, invalid room definition (%s)\n", currentLine)
				os.Exit(1)
			}
			colony.CreateRoom(roomName, "end", xPos, yPos)
			continue
		}
		// Process room connections
		if strings.Contains(currentLine, "-") && strings.Count(currentLine, "-") == 1 {
			connectionParts := strings.Split(currentLine, "-")
			if !colony.ConnectRooms(connectionParts[0], connectionParts[1]) {
				os.Exit(1)
			}
			continue
		}

		// Process standard room definitions
		roomComponents := strings.Fields(currentLine)
		if len(roomComponents) != 3 && (!isStart || !isEnd) {
			// Skip lines that dont match room definitions or links. (handles unknown commands)
			continue
		}
		roomName := roomComponents[0]
		xPos, _ := strconv.Atoi(roomComponents[1])
		yPos, _ := strconv.Atoi(roomComponents[2])
		colony.CreateRoom(roomName, "normal", xPos, yPos)
	}

	if totalAnts == 0 {
		fmt.Println("Data format issue - No ants specified!")
		os.Exit(1)
	}

	if !startRoomDefined || !endRoomDefined {
		fmt.Println("ERROR: invalid data format, no ##start/##end room found")
		os.Exit(1)
	}

	colony.SetupAnts(totalAnts)
	return true, fileContent
}
