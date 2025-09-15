package funcs

import (
	"fmt"
	"math"
	"os"
)

// Add room to end of search queue
func (queue *SearchQueue) AddToQueue(room *Room) {
	queue.roomList = append(queue.roomList, room)
}

// Remove and return first room from queue
func (queue *SearchQueue) RemoveFromQueue() *Room {
	if len(queue.roomList) == 0 {
		return nil
	}
	room := queue.roomList[0]
	queue.roomList = queue.roomList[1:]
	return room
}

// Compute shortest paths from endRoom to all reachable rooms
func (colony *Colony) FindBestPath() {
	queue := &SearchQueue{}
	processedRooms := make(map[*Room]bool)
	foundRooms := make(map[*Room]bool)

	// Begin BFS from the end room with distance 0
	queue.AddToQueue(colony.endRoom)
	colony.roomPaths[colony.endRoom] = 0 // end room should be distance 0
	foundRooms[colony.endRoom] = true

	// Ensure start room is accessible
	if !processedRooms[colony.startRoom] {
		fmt.Println("ERROR: No path exists from start to end")
		os.Exit(0)
	}

	for {
		current := queue.RemoveFromQueue()
		if current == nil {
			break
		}

		tunnel := current.tunnels.firstNode
		for tunnel != nil {
			neighbor := tunnel.data
			newDistance := colony.roomPaths[current] + 1

			if !foundRooms[neighbor] {
				colony.roomPaths[neighbor] = newDistance
				foundRooms[neighbor] = true
				if !neighbor.isStart {
					queue.AddToQueue(neighbor)
				}
			} else if newDistance < colony.roomPaths[neighbor] {
				colony.roomPaths[neighbor] = newDistance
			}

			if neighbor == colony.startRoom {
				processedRooms[neighbor] = true
			}
			tunnel = tunnel.nextConnection
		}
		processedRooms[current] = true
	}

}

// Ckeck if all ants have reached the destination
func (colony *Colony) CheckAllAntsAtDestination() bool {
	for _, ant := range colony.ants {
		if ant.currentRoom != colony.endRoom {
			return false
		}
	}
	return true
}

// Clear all tunnel access restrictions
func (colony *Colony) ClearTunnelRestrictions() {
	for _, room := range colony.rooms {
		for tunnelName := range room.accessMap {
			room.accessMap[tunnelName] = false
		}
	}
}

// Reset all ants to starting position
func (colony *Colony) ResetAllAnts() {
	for i := 0; i < colony.antCount; i++ {
		colony.ants[i].currentRoom = colony.startRoom
		// Clear movement history
		for room := range colony.ants[i].visitedRoom {
			colony.ants[i].visitedRoom[room] = false
		}
		colony.ants[i].visitedRoom[colony.startRoom] = false
		colony.ants[i].inMotion = true
		colony.ants[i].hasCompletedMove = false
	}
	colony.ClearTunnelRestrictions()
}

// Initialize path distances for pathfinding algorithm
func (colony *Colony) InitializePathDistance() {
	for _, room := range colony.rooms {
		colony.roomPaths[room] = math.MaxInt32
	}
}

// Prepare colony for pathfinding simulation
func (colony *Colony) startPathFinding() {
	colony.InitializePathDistance()
	colony.FindBestPath()
}

// Verify colony meets requirements for simulation
func (colony *Colony) VerifyColonyIntegrity() bool {
	// Check if start and end rooms exist
	if colony.startRoom == nil || colony.endRoom == nil {
		fmt.Println("ERROR: Missing start or end room")
		return false
	}

	// Validate coordinate are uniqe
	if !colony.CheckCoordinateUniqueness() {
		fmt.Println("ERROR: Duplicate coordinates found")
		return false
	}

	// Check if ants are present
	if colony.antCount <= 0 {
		fmt.Println("ERROR: No ants to move")
		return false
	}

	return true
}

// verify if pathfinding did work correctly
func (colony *Colony) TestPathfinding() {
	if !colony.VerifyColonyIntegrity() {
		fmt.Println("ERROR: Invalid colony structure")
		return
	}

	colony.startPathFinding()

	// Display path distances for verification
	fmt.Printf("\nPath distances calculated:\n")
	fmt.Printf("Start room distance: %d\n", colony.roomPaths[colony.startRoom])
	fmt.Printf("End room distance: %d\n", colony.roomPaths[colony.endRoom])

	if colony.roomPaths[colony.startRoom] != math.MaxInt32 {
		fmt.Printf("Path from start to end exists! Length: %d steps\n", colony.roomPaths[colony.startRoom])
	} else {
		fmt.Println("No path found from start to end")
	}
}
