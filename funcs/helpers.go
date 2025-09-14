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
		fmt.Println("ERROR: invalid data format, no start/end room found")
		os.Exit(1)
	}

	colony.SetupAnts(totalAnts)
	return true, fileContent
}

func (colony *Colony) SetupAnts(antCount int) {
	colony.antCount = antCount
	colony.ants = make([]*Ant, colony.antCount)

	// Initializing each ant
	for i := range antCount {
		colony.ants[i] = new(Ant)
		colony.ants[i].currentRoom = colony.startRoom
		colony.ants[i].visitedRoom = make(map[*Room]bool)
		colony.ants[i].visitedRoom[colony.startRoom] = false
		colony.ants[i].inMotion = true
		colony.ants[i].antNumber = i + 1
	}
}

func (colony *Colony) Setup() {
	colony.rooms = make(map[string]*Room)
	colony.roomPaths = make(map[*Room]int)
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

func (farm *Colony) ConnectRooms(room1Name string, room2Name string) bool {
	room1 := farm.rooms[room1Name]
	room2 := farm.rooms[room2Name]

	if room1 == nil || room2 == nil {
		fmt.Println("ERROR: invalid data format, invalid room definition")
		return false
	}

	if room1.HasConnectionTo(room2) {
		fmt.Printf("ERROR: duplicate connection between rooms '%s' and '%s'\n", room1Name, room2Name)
		return false
	}

	room1.tunnels.AppendRoom(room2)
	room2.tunnels.AppendRoom(room1)

	return true
}

// checks if a room has a tunnel to the target room
func (room *Room) HasConnectionTo(target *Room) bool {
	currentTunnel := room.tunnels.firstNode

	for currentTunnel != nil {
		if currentTunnel.data.roomName == target.roomName {
			return true
		}
		currentTunnel = currentTunnel.nextConnection
	}
	return false
}

func (list *TunnelList) AppendRoom(room *Room) {
	newTunnel := &Tunnel{data: room, nextConnection: nil}

	if list.firstNode == nil {
		list.firstNode = newTunnel
		return
	}

	current := list.firstNode
	for current.nextConnection != nil {
		current = current.nextConnection
	}
	current.nextConnection = newTunnel
}
