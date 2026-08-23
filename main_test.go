package main

import (
	"testing"

	"lem-in/parser"
	"lem-in/simulator"
	"lem-in/solver"
	"lem-in/verify"
)

var expectedTurns = map[string]int{
	"testdata/test0.txt":            4,
	"testdata/test1.txt":            5,
	"testdata/test2.txt":            3,
	"testdata/test3_branching.txt":  8,
	"testdata/test4_bottleneck.txt": 6,
	"testdata/example02.txt":        11,
	"testdata/example05.txt":        8,
}

func TestPipeline_AllFixtures(t *testing.T) {
	for path, wantTurns := range expectedTurns {
		path, wantTurns := path, wantTurns
		t.Run(path, func(t *testing.T) {
			colony, err := parser.ParseFile(path)
			if err != nil {
				t.Fatalf("ParseFile: %v", err)
			}

			assignment := solver.Solve(colony.Links, colony.Start, colony.End, colony.NumAnts)
			if assignment == nil {
				t.Fatalf("Solve returned nil (no path found) for a fixture that should have one)")
			}

			turns := simulator.Run(assignment, colony.NumAnts)

			if err := verify.Check(colony, colony.NumAnts, turns); err != nil {
				t.Errorf("output failed replay verification: %v", err)
			}

			if len(turns) != wantTurns {
				t.Errorf("took %d turns, want %d", len(turns), wantTurns)
			}
		})
	}
}

func TestPipeline_RejectsInvalidInput(t *testing.T) {
	if _, err := parser.ParseFile("testdata/does_not_exist.txt"); err == nil {
		t.Error("expected an error for a missing file, got nil")
	}
}
