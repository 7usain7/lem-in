## Lem-in: Ant Farm Pathfinding Algorithm

A Go implementation of a digital ant farm that finds the optimal path for multiple ants to navigate from a start room to an end room through a network of connected tunnels.

## Key Features

- BFS Pathfinding: Uses Breadth-First Search to calculate shortest distances from end to all rooms.

- Traffic Management: Prevents ants from colliding or blocking each other.

- Dual Algorithm Comparison: Runs both basic and optimized algorithms, selecting the best result.

- Input Validation: Comprehensive error checking for invalid data formats.

- Flexible Room Networks: Supports complex tunnel systems with multiple paths.

## How It Works

1. Data Structures

**The program uses several key data structures:**

- Colony: Contains rooms, ants, and path information
- Room: Represents a location with coordinates and tunnel connections
- Ant: Tracks ant position and movement history
- Tunnel: Links between rooms (linked lists)
- SearchQueue: Queue for BFS pathfinding

2. Algorithm and validation Steps

- Parse Input: Read and validate the input file format
- Validate Colony Structure
- Build Network: Create rooms and establish tunnel connections
- Calculate Paths: Use BFS to find shortest distances from end room to all other rooms
- Move Ants: Simulate ant movement using two different strategies
- Compare Results: Select the approach that requires fewer steps

3. Movement Rules

- Only one ant per room (except start/end rooms)
- Each tunnel can only be used once **per turn**
- Ants cannot revisit rooms they already been in
- Ants move toward rooms with shorter distances to the destination

## How to run

` go run main.go <file_Name>`

## Team Members

- Ali Hussain #alimadan
- Hussain Abdulrasool #habdulras
