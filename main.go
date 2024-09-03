package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Room struct {
	Name string
	X    int
	Y    int
}
type Farm struct {
	Rooms map[string]Room
	Links map[string][]string
}
type Path struct {
	rooms    []string
	roomsNum int
}

var (
	startCounter int
	endCounter   int
	startRoom    Room
	endRoom      Room
	coords       = make(map[[2]int]bool)
	farm         Farm
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	inputFile := os.Args[1]
	lines, err := ReadFile(inputFile)
	if err != nil {
		fmt.Println("error while opening the file")
		os.Exit(0)
	}
	farm = Farm{
		Rooms: make(map[string]Room),
		Links: make(map[string][]string),
	}

	if err := HandelInput(lines); err != nil {
		fmt.Println("Error: ", err)
	}

	if err := HandleRoomInfo(lines); err != nil {
		fmt.Println("Error: ", err)
	}

}

func ReadFile(FileName string) ([]string, error) {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func HandelInput(lines []string) error {
	var antCount string
	printLine := false

	for _, line := range lines {
		// Ignore lines containing '$'
		if strings.Contains(line, "$") {
			continue
		}

		// Ignore lines starting with 'L' followed by a number
		if strings.HasPrefix(line, "L") && len(line) > 1 && unicode.IsDigit(rune(line[1])) {
			continue
		}

		// Ignore lines that start with a single '#'
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}

		// Check for ant count before '##start'
		if antCount == "" && line != "" && line != "##start" {
			antCount = line
			fmt.Println(antCount)
			continue
		}

		if line == "##start" {
			if antCount == "" {
				return fmt.Errorf("no ants found")

			}
			// Print '##start' line
			fmt.Println(line)
			printLine = true
			continue
		}

		if printLine {
			fmt.Println(line)
		}

		// Stop processing if a line starts with 'L'
		if strings.HasPrefix(line, "L") {
			break
		}
	}
	return nil
}

func HandleRoomInfo(lines []string) error {
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "##start" {
			if startCounter > 0 {
				return fmt.Errorf("multiple start rooms found")
			}
			if i >= len(lines)-1 {
				return fmt.Errorf("invalid syntax")
			}
			i++
			startRoomInfo := strings.TrimSpace(lines[i])
			StartInfo := strings.Split(startRoomInfo, " ")
			if len(StartInfo) != 3 {
				return fmt.Errorf("invalid start info")
			}
			err := addRoom(StartInfo[0], StartInfo[1], StartInfo[2], "start")
			if err != nil {
				return err
			}
			startCounter++
			continue // Skip the next iteration
		}

		if line == "##end" {
			if endCounter > 0 {
				return fmt.Errorf("multiple end rooms found")
			}
			if i >= len(lines)-1 {
				return fmt.Errorf("invalid syntax")
			}
			i++
			endRoomInfo := strings.TrimSpace(lines[i])
			endInfo := strings.Split(endRoomInfo, " ")
			if len(endInfo) != 3 {
				return fmt.Errorf("invalid end info")
			}
			err := addRoom(endInfo[0], endInfo[1], endInfo[2], "end")
			if err != nil {
				return err
			}
			endCounter++
			continue // Skip the next iteration
		}

		// Process normal rooms (not start or end)
		if strings.Contains(line, " ") && !strings.Contains(line, "-") && line != "##start" && line != "##end" {
			NormalRoomInfo := strings.Split(strings.TrimSpace(line), " ")
			if len(NormalRoomInfo) != 3 {
				return fmt.Errorf("invalid room info")
			}
			err := addRoom(NormalRoomInfo[0], NormalRoomInfo[1], NormalRoomInfo[2], "Normal")
			if err != nil {
				return err
			}
		}
		if strings.Contains(line, "-") {
			err := parseLink(line)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addRoom(name string, xi string, yi string, roomType string) error {
	x, err := strconv.Atoi(xi)
	if err != nil {
		return fmt.Errorf("invalid x coordinate for room %s", name)
	}
	y, err := strconv.Atoi(yi)
	if err != nil {
		return fmt.Errorf("invalid y coordinate for room %s", name)
	}
	if _, exists := farm.Rooms[name]; exists {
		fmt.Println(farm.Rooms)
		fmt.Println(name, x, y)
		return fmt.Errorf("duplicate room name: %s", name)
	}
	coord := [2]int{x, y}
	if coords[coord] {
		return fmt.Errorf("duplicate room coordinates: %d %d", x, y)
	}
	room := Room{Name: name, X: x, Y: y}
	farm.Rooms[name] = room
	coords[coord] = true
	//fmt.Println(room.Name)
	if roomType == "start" {
		startRoom = room
	}
	if roomType == "end" {
		endRoom = room
	}
	return nil
}

// parseLink processes a link definition line
func parseLink(line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid link format: %s", line)
	}
	room1, room2 := parts[0], parts[1]
	if _, exists := farm.Rooms[room1]; !exists {
		return fmt.Errorf("unknown room in link: %s", room1)
	}
	if _, exists := farm.Rooms[room2]; !exists {
		return fmt.Errorf("unknown room in link: %s", room2)
	}
	farm.Links[room1] = append(farm.Links[room1], room2)
	farm.Links[room2] = append(farm.Links[room2], room1)
	return nil
}

// reconstructPath builds the path from start to end using the parent map
func reconstructPath(parent map[string]string, start, end string) []string {
    path := []string{end}
    for current := end; current != start; {
        current = parent[current]
        path = append([]string{current}, path...)
    }
    return path
}