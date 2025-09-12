package funcs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseInput(filename string) (*Colony, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	colony := &Colony{}

	scanner := bufio.NewScanner(file)
	lineIndex := 0
	nextRoomIsStart := false
	nextRoomIsEnd := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		colony.Input = append(colony.Input, line)

		if line == "" {
			continue
		}

		if line == "##start" {
			nextRoomIsStart = true
			continue
		}
		if line == "##end" {
			nextRoomIsEnd = true
			continue
		}

		// Skip other comments if exist
		if strings.HasPrefix(line, "#") {
			continue
		}

		// First line must be the number of ants
		if lineIndex == 0 {
			AntsNum, err := strconv.Atoi(line)
			if err != nil || AntsNum <= 0 {
				return nil, fmt.Errorf("ERROR: Invalid number of ants: %s", line)
			}
			colony.AntsNum = AntsNum
			lineIndex++
			continue
		}
	}

	// Errors returns
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// if colony.Start == "" {
	// 	return nil, fmt.Errorf("ERROR: No start room found")
	// }

	// if colony.End == "" {
	// 	return nil, fmt.Errorf("ERROR: No end room found")
	// }

	// if colony.Start == colony.End {
	// 	return nil, fmt.Errorf("ERROR: Invalid input: start and end rooms are the same")
	// }

	fmt.Println(colony.AntsNum)
	fmt.Println(nextRoomIsStart)
	fmt.Println(nextRoomIsEnd)

	fmt.Println("===========================================")

	return colony, nil
}

func DisplayColony(colony *Colony) {
	for _, line := range colony.Input {
		fmt.Println(line)
	}
	fmt.Println()
}
