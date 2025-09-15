package funcs

import (
	"fmt"
	"math"
	"os"
	"sort"
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

	// Validate if coordinate are unique
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

// Find the best next room for an ant to move toward
func (colony *Colony) FindOptimalNextRoom(ant *Ant) *Room {
	minDistance := math.MaxInt32
	bestRoom := ant.currentRoom.tunnels.firstNode.data

	tunnel := ant.currentRoom.tunnels.firstNode
	for tunnel != nil {
		neighbor := tunnel.data
		if colony.roomPaths[neighbor] <= minDistance &&
			!ant.visitedRoom[neighbor] {
			bestRoom = neighbor
			minDistance = colony.roomPaths[bestRoom]
		}
		tunnel = tunnel.nextConnection
	}
	return bestRoom
}

// Validate if ant is allowed to move to specified room
func validateAntMovement(ant *Ant, targetRoom *Room) bool {
	// Ant can move to end room
	if targetRoom.isEnd {
		return ant.inMotion && !ant.hasCompletedMove
	}

	// For non-end rooms, check all conditions
	return targetRoom.isUnoccupied &&
		!ant.visitedRoom[targetRoom] &&
		!targetRoom.isStart &&
		!ant.currentRoom.accessMap[targetRoom.roomName] &&
		ant.inMotion &&
		!ant.hasCompletedMove
}

// Execute basic ant movement through the colony system
func (colony *Colony) StartAntMovment() {
	for _, ant := range colony.ants {
		if !ant.inMotion || ant.hasCompletedMove {
			continue
		}

		nextTarget := colony.FindOptimalNextRoom(ant)

		// Perform movement if conditions are met
		if validateAntMovement(ant, nextTarget) {
			// Block tunnel to prevent conflicts
			ant.currentRoom.accessMap[nextTarget.roomName] = true
			ant.visitedRoom[ant.currentRoom] = true
			ant.currentRoom.isUnoccupied = true
			ant.currentRoom = nextTarget
			nextTarget.isUnoccupied = false

			// Update ant status if destination reached
			if ant.currentRoom.isEnd {
				ant.inMotion = false
			}
			ant.hasCompletedMove = true
		}
	}

	// Clear access restrictions and recalculate paths
	colony.ClearTunnelRestrictions()
	colony.FindBestPath()
}

// Generate formatted string of current ant locations
func (colony *Colony) GenerateMovementOutput() string {
	outputString := ""

	// Create numerically sorted list of ants
	orderedAnts := make([]*Ant, colony.antCount)
	copy(orderedAnts, colony.ants)

	// Sort ants by their identification numbers
	sort.SliceStable(orderedAnts, func(first, second int) bool {
		return orderedAnts[first].antNumber < orderedAnts[second].antNumber
	})

	// Format output for ants that completed movement
	for _, ant := range orderedAnts {
		if ant.hasCompletedMove {
			outputString += fmt.Sprintf("L%d-%s ", ant.antNumber, ant.currentRoom.roomName)
			ant.hasCompletedMove = false
		}
	}
	return outputString + "\n"
}

// verify if pathfinding did work correctly
func (colony *Colony) TestPathfinding() {
	if !colony.VerifyColonyIntegrity() {
		fmt.Println("ERROR: Invalid colony structure")
		return
	}

	colony.startPathFinding()

	steps := 0
	for !colony.CheckAllAntsAtDestination() {
		steps++
		colony.StartAntMovment()
		fmt.Print(colony.GenerateMovementOutput())
	}
	fmt.Printf("Steps taken: %d\n", steps)
}
