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

	// Check accessibility AFTER BFS completes
	if !processedRooms[colony.startRoom] {
		fmt.Println("ERROR: No path exists from start to end")
		os.Exit(0)
	}
}

// Find the best next room for an ant to move toward
func (colony *Colony) FindOptimalNextRoom(ant *Ant, strictMode bool) *Room {
	minDistance := math.MaxInt32
	bestRoom := ant.currentRoom.tunnels.firstNode.data

	tunnel := ant.currentRoom.tunnels.firstNode
	for tunnel != nil {
		neighbor := tunnel.data
		if strictMode {
			if colony.roomPaths[neighbor] < minDistance &&
				!ant.visitedRoom[neighbor] {
				bestRoom = neighbor
				minDistance = colony.roomPaths[bestRoom]
			}
		} else {
			if colony.roomPaths[neighbor] <= minDistance &&
				!ant.visitedRoom[neighbor] {
				bestRoom = neighbor
				minDistance = colony.roomPaths[bestRoom]
			}
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

var movementHappened = false

// Execute ant movement through the colony system
func (colony *Colony) StartAntMovment(optimizationMode bool) {
	workingAnts := make([]*Ant, colony.antCount)
	copy(workingAnts, colony.ants)

	// Sort ants by tunnel availability at current location
	sort.SliceStable(workingAnts, func(i, j int) bool {
		return colony.CountProgressiveTunnels(workingAnts[i].currentRoom) <
			colony.CountProgressiveTunnels(workingAnts[j].currentRoom)
	})

	for antIdx := 0; antIdx < len(workingAnts); antIdx++ {
		nextTarget := colony.FindOptimalNextRoom(workingAnts[antIdx], optimizationMode)

		// Adjust path weights for traffic management
		if nextTarget != colony.endRoom || workingAnts[antIdx].currentRoom.isStart {
			colony.roomPaths[nextTarget]++
		}

		// Perform movement if conditions are met
		if validateAntMovement(workingAnts[antIdx], nextTarget) {
			// Block tunnel to prevent conflicts
			workingAnts[antIdx].currentRoom.accessMap[nextTarget.roomName] = true
			workingAnts[antIdx].visitedRoom[workingAnts[antIdx].currentRoom] = true
			workingAnts[antIdx].currentRoom.isUnoccupied = true
			workingAnts[antIdx].currentRoom = nextTarget
			nextTarget.isUnoccupied = false

			// Update ant status if destination reached
			if workingAnts[antIdx].currentRoom.isEnd {
				workingAnts[antIdx].inMotion = false
			}
			workingAnts[antIdx].hasCompletedMove = true
			movementHappened = true
		}

		// Restart loop if movement detected and at final ant
		if antIdx == len(workingAnts)-1 && movementHappened {
			antIdx = 0
			movementHappened = false
		}
	}

	// Clear access restrictions and recalculate paths
	colony.ClearTunnelRestrictions()
	colony.FindBestPath()
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

// Verify all ants have reached the destination
func (colony *Colony) CheckAllAntsAtDestination() bool {
	for _, ant := range colony.ants {
		if ant.currentRoom != colony.endRoom {
			return false
		}
	}
	return true
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

// Execute complete pathfinding simulation with dual mode comparison
func (colony *Colony) RunOptimizedPathfinding() {
	basicSteps, advancedSteps := 0, 0
	basicPath, advancedPath := "", ""
	stepCounter := 0

	// Execute basic pathfinding approach
	for !colony.CheckAllAntsAtDestination() {
		basicSteps++
		colony.StartAntMovment(false)
		basicPath += colony.GenerateMovementOutput()
	}

	// Reset and execute optimized approach
	colony.ResetAllAnts()
	for !colony.CheckAllAntsAtDestination() {
		stepCounter++
		advancedSteps++
		if advancedSteps > basicSteps {
			break
		}
		colony.StartAntMovment(true)
		advancedPath += colony.GenerateMovementOutput()
	}

	// Display final result
	if basicSteps == advancedSteps {
		fmt.Printf("\n%s\nSteps taken: %d\n", advancedPath, advancedSteps)
	} else if basicSteps < advancedSteps {
		fmt.Printf("\n%s\nSteps taken: %d\n", basicPath, basicSteps)
	} else {
		fmt.Printf("\n%s\nSteps taken: %d\n", advancedPath, advancedSteps)
	}
}

// Prepare colony for pathfinding simulation
func (colony *Colony) PrepareSimulation() {
	colony.InitializePathDistances()
	colony.FindBestPath()
}

// Verify colony meets requirements for simulation
func (colony *Colony) VerifyColonyIntegrity() bool {
	// Confirm start and end rooms exist
	if colony.startRoom == nil || colony.endRoom == nil {
		return false
	}

	// Validate coordinate uniqueness
	if !colony.CheckCoordinateUniqueness() {
		return false
	}

	// Confirm ants are present
	if colony.antCount <= 0 {
		return false
	}

	return true
}

// Execute pathfinding
func (colony *Colony) LaunchPathfindingSimulation() {
	if !colony.VerifyColonyIntegrity() {
		fmt.Println("ERROR: Invalid colony structure")
		return
	}

	colony.PrepareSimulation()
	colony.RunOptimizedPathfinding()
}
