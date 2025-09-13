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
	}
}

func (colony *Colony) CreateRoom(name string, roomType string, x int, y int) {
	if _, exists := colony.rooms[name]; exists {
		fmt.Printf("ERROR: invalid data format, duplicate room name '%s'\n", name)
		os.Exit(1)
	}

	newRoom := &Room{
		roomName:     name,
		xCoord:       x,
		yCoord:       y,
		tunnels:      &TunnelList{},
		accessMap:    make(map[string]bool),
		isUnoccupied: true,
	}

	switch roomType {
	case "start":
		newRoom.isStart = true
		newRoom.isEnd = false
		colony.rooms[name] = newRoom
		colony.startRoom = newRoom

	case "end":
		newRoom.isStart = false
		newRoom.isEnd = true
		colony.rooms[name] = newRoom
		colony.endRoom = newRoom

	case "normal":
		newRoom.isStart = false
		newRoom.isEnd = false
		colony.rooms[name] = newRoom

	default:
		fmt.Printf("ERROR: invalid data format, invalid room type '%s'\n", roomType)
		os.Exit(1)
	}
}
