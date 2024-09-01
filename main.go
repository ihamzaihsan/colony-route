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

	if err := HandelRoomInfo(lines); err != nil {
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

func HandelRoomInfo(lines []string) error {
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "##start" {
			if startCounter > 0 {
				return  fmt.Errorf("multiple start rooms found")
			}
			if i >= len(lines)-3 {
				return fmt.Errorf("invalid sentyx")
			}
			startRoomInfo := strings.TrimSpace(lines[i+1])
			if len(startRoomInfo) != 5 {
				return fmt.Errorf("invalid start info")
			}
			err:= addRoom(string(startRoomInfo[0]), string(startRoomInfo[2]), string(startRoomInfo[4]))
			if err != nil{
				return err
			}
			

		}
	}
	return nil
}
func addRoom(name string, xi string, yi string) error {
	x, err := strconv.Atoi(xi)
	if err != nil {
		return fmt.Errorf("invalid x coordinate for room %s", name)
	}
	y, err := strconv.Atoi(yi)
	if err != nil {
		return fmt.Errorf("invalid y coordinate for room %s", name)
	}
	if _, exists := farm.Rooms[name]; exists {
		return fmt.Errorf("duplicate room name: %s", name)
	}
	room := Room{Name: name, X: x, Y: y}
	farm.Rooms[name] = room
	fmt.Println(room.Name)
	return nil
}
