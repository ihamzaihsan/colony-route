package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	inputFile := os.Args[1]
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var antCount string
	printLine := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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
				fmt.Println("Error: No ants found")
				return
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

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
