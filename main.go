package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/parser"
	"lem-in/solver"
	"lem-in/simulator"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format, usage: lem-in <filename>")
		os.Exit(1)
	}

	colony, err := parser.ParseFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assignment := solver.Solve(colony.Links, colony.Start, colony.End, colony.NumAnts)
	if assignment == nil {
		fmt.Println("ERROR: invalid data format, no path between ##start and ##end")
		os.Exit(1)
	}

	printFileContent(colony)
	fmt.Println()

	for _, line := range simulator.Run(assignment, colony.NumAnts) {
		fmt.Println(line)
	}
}

func printFileContent(c *parser.Colony) {
	skippedAntCount := false
	for _, raw := range c.RawLines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if !skippedAntCount {
			if _, err := strconv.Atoi(line); err == nil {
				fmt.Println(line)
				skippedAntCount = true
				continue
			}
		}
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}
		fmt.Println(line)
	}
}
