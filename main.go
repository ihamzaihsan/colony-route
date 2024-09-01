package main

import (
	"fmt"
	"io/ioutil"
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
	lines, err := ReadFile(inputFile)
	if err != nil {
		fmt.Println("error while opening the file")
		os.Exit(0)
	}

	if err := HandelInput(lines); err != nil{
		fmt.Println("Error: ", err )
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
