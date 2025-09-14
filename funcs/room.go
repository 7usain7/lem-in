package funcs

import (
	"fmt"
	"os"
)

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

// Append new room connection to tunnel list
func (tunnelList *TunnelList) AppendConnection(roomToAdd *Room) {
	newTunnel := &Tunnel{
		data:           roomToAdd,
		nextConnection: nil,
	}

	if tunnelList.firstNode == nil {
		tunnelList.firstNode = newTunnel
		return
	}

	// Navigate to end of list and append
	current := tunnelList.firstNode
	for current.nextConnection != nil {
		current = current.nextConnection
	}
	current.nextConnection = newTunnel
}

// Construct and register new room in colony
func (colony *Colony) BuildRoom(roomName string, roomType string, xCoord int, yCoord int) {
	// Prevent duplicate room names
	if _, exists := colony.rooms[roomName]; exists {
		fmt.Printf("ERROR: invalid data format, duplicate room name '%s'\n", roomName)
		os.Exit(1)
	}

	// Construct room with specified attributes
	newRoom := &Room{
		roomName:     roomName,
		xCoord:       xCoord,
		yCoord:       yCoord,
		tunnels:      &TunnelList{},
		accessMap:    make(map[string]bool),
		isUnoccupied: true,
	}

	// Configure room based on type
	switch roomType {
	case "start":
		newRoom.isStart = true
		newRoom.isEnd = false
		colony.rooms[roomName] = newRoom
		colony.startRoom = newRoom

	case "end":
		newRoom.isStart = false
		newRoom.isEnd = true
		colony.rooms[roomName] = newRoom
		colony.endRoom = newRoom

	case "normal":
		newRoom.isStart = false
		newRoom.isEnd = false
		colony.rooms[roomName] = newRoom

	default:
		fmt.Printf("ERROR: invalid data format, invalid room type '%s'\n", roomType)
		os.Exit(1)
	}
}

// Establish bidirectional tunnel between two rooms
func (colony *Colony) LinkRooms(fromRoomName string, toRoomName string) bool {
	fromRoom := colony.rooms[fromRoomName]
	toRoom := colony.rooms[toRoomName]

	// Verify both rooms are valid
	if fromRoom == nil || toRoom == nil {
		fmt.Printf("ERROR: invalid data format. Cannot link rooms '%s' and '%s' because one or both are not defined.\n", fromRoomName, toRoomName)
		return false
	}

	// Check for duplicate connections
	if fromRoom.HasConnectionTo(toRoom) {
		fmt.Printf("ERROR: duplicate connection between rooms '%s' and '%s'\n",
			fromRoomName, toRoomName)
		return false
	}

	// Create bidirectional tunnel connection
	fromRoom.tunnels.AppendConnection(toRoom)
	toRoom.tunnels.AppendConnection(fromRoom)

	return true
}
