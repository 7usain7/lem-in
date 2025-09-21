package funcs

import (
	"fmt"
	"os"
	"strings"
)

func (colony *Colony) CreateRoom(name string, roomType string, x int, y int) {
	// Validate room name format
	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") || strings.Contains(name, " ") {
		fmt.Printf("ERROR: invalid data format, invalid room name '%s'\n", name)
		os.Exit(1)
	}

	if _, exists := colony.rooms[name]; exists {
		fmt.Printf("ERROR: invalid data format, duplicate room name '%s'\n", name)
		os.Exit(1)
	}

	// Check for duplicate coordinates
	for _, room := range colony.rooms {
		if room.xCoord == x && room.yCoord == y {
			fmt.Printf("ERROR: invalid data format, duplicate coordinates (%d, %d)\n", x, y)
			os.Exit(1)
		}
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

func (colony *Colony) ConnectRooms(room1Name string, room2Name string) bool {
	room1 := colony.rooms[room1Name]
	room2 := colony.rooms[room2Name]

	if room1 == nil || room2 == nil {
		fmt.Printf("ERROR: invalid data format. \n")
		return false
	}

	if room1.HasConnectionTo(room2) {
		fmt.Printf("ERROR: invalid data format. Duplicate connection between rooms")
		return false
	}

	room1.tunnels.AppendRoom(room2)
	room2.tunnels.AppendRoom(room1)

	return true
}
