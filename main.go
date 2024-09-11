package main

import (
	"container/list"
	"fmt"
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
	antCount     string
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
		fmt.Println("Error: error while opening the file")
		return
	}

	// Remove empty lines
	lines = removeEmptyLines(lines)

	farm = Farm{
		Rooms: make(map[string]Room),
		Links: make(map[string][]string),
	}
	if !isStartEndExist(lines) {
		fmt.Println("no start or end room is found")
		os.Exit(0)
	}
	num := strings.TrimSpace(lines[0])
	numAnts, err := strconv.Atoi(num)

	if err != nil || numAnts <= 0 || numAnts > 10000 {
		fmt.Println("Error: Invalid number of ants")
		return
	}
	if err := HandelInput(lines); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if err := HandleRoomInfo(lines); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if startCounter != 1 {
		fmt.Println("Error: only one start room is valid")
		os.Exit(0)
	}
	if endCounter != 1 {
		fmt.Println("Error: only one end room is valid")
		os.Exit(0)
	}
	paths := findPathsEdmondsKarp(startRoom.Name, endRoom.Name)
	if len(paths) == 0 {
		fmt.Println("Error: No valid paths found")
		return
	}
	PrintData(lines)
	totalSteps := moveAnts(numAnts, paths)
	fmt.Printf("\nTotal steps: %d\n", totalSteps)
}

func ReadFile(FileName string) ([]string, error) {
	data, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func HandelInput(lines []string) error {
	for _, line := range lines {
		// Check for ant count before '##start'
		if antCount == "" && line != "##start" {
			antCount = line
			// fmt.Println(antCount)
			continue
		}
		if line == "##start" {
			if antCount == "" {
				return fmt.Errorf("no ants found")
			}
		}
		// Stop processing if a line starts with 'L'
		if strings.HasPrefix(line, "L") {
			return fmt.Errorf("wrong name: L")
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
		return fmt.Errorf("duplicate room name: %s", name)
	}
	coord := [2]int{x, y}
	if coords[coord] {
		return fmt.Errorf("duplicate room coordinates: %d %d", x, y)
	}
	room := Room{Name: name, X: x, Y: y}
	farm.Rooms[name] = room
	coords[coord] = true
	// fmt.Println(room.Name)
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

	if room1 == room2 {
		return fmt.Errorf("self-loop detected in link: %s", line)
	}

	if _, exists := farm.Rooms[room1]; !exists {
		return fmt.Errorf("unknown room in link: %s", room1)
	}
	if _, exists := farm.Rooms[room2]; !exists {
		return fmt.Errorf("unknown room in link: %s", room2)
	}

	// Check for duplicate connections
	for _, connectedRoom := range farm.Links[room1] {
		if connectedRoom == room2 {
			return fmt.Errorf("duplicate link detected: %s", line)
		}
	}
	for _, connectedRoom := range farm.Links[room2] {
		if connectedRoom == room1 {
			return fmt.Errorf("duplicate link detected: %s", line)
		}
	}

	// Add the connection to the map
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

// bfs performs a breadth-first search to find an augmenting path
func bfs(start, end string, graph map[string]map[string]int) ([]string, bool) {
	queue := list.New()
	queue.PushBack(start)
	visited := make(map[string]bool)
	parent := make(map[string]string)
	visited[start] = true
	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(string)
		if current == end {
			return reconstructPath(parent, start, end), true
		}
		for neighbor, capacity := range graph[current] {
			if !visited[neighbor] && capacity > 0 {
				visited[neighbor] = true
				parent[neighbor] = current
				queue.PushBack(neighbor)
			}
		}
	}
	return nil, false
}

// findPathsEdmondsKarp implements the Edmonds-Karp algorithm to find multiple paths
func findPathsEdmondsKarp(startRoomName, endRoomName string) []Path {
	var paths []Path
	residualGraph := make(map[string]map[string]int)

	// Initialize residual graph
	for room, links := range farm.Links {
		residualGraph[room] = make(map[string]int)
		for _, link := range links {
			residualGraph[room][link] = 1
		}
	}
	// Find augmenting paths until no more are found
	for {
		path, found := bfs(startRoomName, endRoomName, residualGraph)
		if !found {
			break
		}
		paths = append(paths, Path{rooms: path, roomsNum: len(path)})

		// Update residual graph
		for i := 0; i < len(path)-1; i++ {
			residualGraph[path[i]][path[i+1]]--
			residualGraph[path[i+1]][path[i]]++
		}
	}

	return paths
}

// moveAnts simulates the movement of ants through the found paths
func moveAnts(numAnts int, paths []Path) int {
	antPath := make(map[int]int)
	antPosition := make(map[int]int)
	InPath := make([]int, len(paths))
	for antID := 1; antID <= numAnts; antID++ {
		pathIndex := 0
		minCost := InPath[0] + paths[0].roomsNum
		for i := 1; i < len(paths); i++ {
			cost := paths[i].roomsNum + InPath[i]
			if minCost > cost {
				minCost = cost
				pathIndex = i
			}
		}
		antPath[antID] = pathIndex
		InPath[pathIndex]++
	}
	antsOutside := make(map[int][]int)
	for i := 0; i < len(paths); i++ {
		antsOutside[i] = make([]int, 0)
	}
	for i := 1; i <= len(antPath); i++ {
		antsOutside[antPath[i]] = append(antsOutside[antPath[i]], i)
	}
	antsInside := make(map[int][]int)
	var antMoving bool
	var output string
	totalSteps := 0
	for step := 1; ; step++ {
		antMoving = false
		for pathIndex := 0; pathIndex < len(paths); pathIndex++ {
			for j := 0; j < len(antsInside[pathIndex]); j++ {
				if antPosition[antsInside[pathIndex][j]] < paths[pathIndex].roomsNum-1 {
					antMoving = true
					antPosition[antsInside[pathIndex][j]]++
					output += fmt.Sprintf("L%d-%s ", antsInside[pathIndex][j], paths[pathIndex].rooms[antPosition[antsInside[pathIndex][j]]])
				}
			}
		}
		for pathIndex := 0; pathIndex < len(paths); pathIndex++ {
			for len(antsOutside[pathIndex]) != 0 {
				if antPosition[antsOutside[pathIndex][0]] < paths[pathIndex].roomsNum-1 {
					antMoving = true
					antPosition[antsOutside[pathIndex][0]]++
					output += fmt.Sprintf("L%d-%s ", antsOutside[pathIndex][0], paths[pathIndex].rooms[antPosition[antsOutside[pathIndex][0]]])
					antID := antsOutside[pathIndex][0]
					antsInside[pathIndex] = append(antsInside[pathIndex], antID)
					antsOutside[pathIndex] = antsOutside[pathIndex][1:]
					break
				}
			}
		}
		if !antMoving {
			totalSteps = step - 1 // Subtract 1 because the last step doesn't move any ants
			break
		}
		fmt.Println(output)
		output = ""
	}
	return totalSteps
}

func isStartEndExist(lines []string) bool {
	start := false
	end := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "##start" {
			start = true
		}
		if line == "##end" {
			end = true
		}

	}
	return start && end
}

func PrintData(lines []string) {
	printLine := false
	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Ignore lines containing '```
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

		// Print '##start' line
		printLine = true

		if printLine {
			fmt.Println(line)
		}
	}
}

func removeEmptyLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		// Trim spaces and check if the line is not empty
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
