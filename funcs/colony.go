package funcs

func (colony *Colony) SetupAnts(antCount int) {
	colony.antCount = antCount
	colony.ants = make([]*Ant, colony.antCount)

	// Initialize each ant with starting conditions
	for i := 0; i < antCount; i++ {
		colony.ants[i] = &Ant{
			antNumber:        i + 1,
			currentRoom:      colony.startRoom,
			visitedRoom:      make(map[*Room]bool),
			inMotion:         true,
			hasCompletedMove: false,
		}
		colony.ants[i].visitedRoom[colony.startRoom] = false
	}
}

func (colony *Colony) Setup() {
	colony.rooms = make(map[string]*Room)
	colony.roomPaths = make(map[*Room]int)
}

// Verify all rooms have unique coordinate locations
func (colony *Colony) CheckCoordinateUniqueness() bool {
	for _, chamber := range colony.rooms {
		for _, otherChamber := range colony.rooms {
			if chamber.roomName != otherChamber.roomName {
				if chamber.xCoord == otherChamber.xCoord &&
					chamber.yCoord == otherChamber.yCoord {
					return false
				}
			}
		}
	}
	return true
}

// Count tunnels that lead closer to destination
func (colony *Colony) CountProgressiveTunnels(chamber *Room) int {
	progressiveCount := 0
	tunnel := chamber.tunnels.firstNode

	for tunnel != nil {
		// Check if connected room has shorter path to destination
		if colony.roomPaths[tunnel.data] < colony.roomPaths[chamber] {
			progressiveCount++
		}
		tunnel = tunnel.nextConnection
	}
	return progressiveCount
}

// TODO: Initialize path distances for pathfinding algorithm
// func (colony *Colony) InitializePathDistances() {
// 	for _, room := range colony.rooms {
// 		colony.roomPaths[room] = math.MaxInt32
// 	}
// }
