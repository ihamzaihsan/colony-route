# lem-in

## Authors
- mohaabdulla
- nahussain
- hussainali2
-hcheema
## Information

lem-in is a program designed to simulate the movement of ants through an ant farm, finding the quickest path from a start room to an end room. It reads input from a file, validates the data, finds all possible paths using Depth-First Search (DFS), filters unique paths, and simulates the movement of ants along these paths.

## Features
- input parsing: Input Parsing: Reads and validates input data describing rooms and tunnels.
- Pathfinding: Uses DFS to find all possible paths from the start room to the end room.
- Path Filtering: Filters unique paths ensuring each room is visited only once per ant.
- Ant Movement Simulation: Simulates the movement of ants along the shortest paths found.
- Error Hadling: Detects and reports errors such as invalid room formats, duplicate rooms, unknown rooms in links, and missing start or end rooms.

## Usage Instructions

- run the program with a test of you choice or the tests in the exampls directory using the following command:

`go run main.go <your_input_file.txt>` or `go run main.go ./examples/test.txt`

- run the program using the bash script which will show the oputput of all the tests with the following command:

`bash test.sh`



