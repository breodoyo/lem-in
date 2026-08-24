package solver

import "testing"

func link(pairs ...[2]string) map[string][]string {
	m := map[string][]string{}
	for _, p := range pairs {
		m[p[0]] = append(m[p[0]], p[1])
		m[p[1]] = append(m[p[1]], p[0])
	}
	return m
}

func turnsUsed(a *Assignment) int {
	counts := make([]int, len(a.Paths))
	max := 0
	for _, p := range a.AntToPath {
		hops := len(a.Paths[p]) - 1
		finish := hops + counts[p]
		counts[p]++
		if finish > max {
			max = finish
		}
	}
	return max
}

func TestSolve_SinglePath(t *testing.T) {
	links := link([2]string{"start", "a"}, [2]string{"a", "end"})
	a := Solve(links, "start", "end", 4)
	if a == nil {
		t.Fatal("expected an assignment, got nil")
	}
	if len(a.AntToPath) != 4 {
		t.Fatalf("len(AntToPath) = %d, want 4", len(a.AntToPath))
	}
	// hops=2, 4 ants: finish times 2,3,4,5 -> 5 turns.
	if got := turnsUsed(a); got != 5 {
		t.Errorf("turnsUsed = %d, want 5", got)
	}
}

func TestSolve_NoPath(t *testing.T) {
	links := link([2]string{"start", "a"}, [2]string{"b", "end"})
	if a := Solve(links, "start", "end", 3); a != nil {
		t.Fatalf("expected nil (no path exists), got %+v", a)
	}
}

func TestSolve_BalancesAcrossUnequalPaths(t *testing.T) {
	links := link(
		[2]string{"start", "end"},
		[2]string{"start", "1"}, [2]string{"1", "2"}, [2]string{"2", "end"},
	)
	a := Solve(links, "start", "end", 20)
	if a == nil {
		t.Fatal("expected an assignment, got nil")
	}
	if got := turnsUsed(a); got != 11 {
		t.Errorf("turnsUsed = %d, want 11 (hand-computed minimum for hops 1 & 3, 20 ants)", got)
	}
}

func TestSolve_BridgeMonopolyGraph(t *testing.T) {
	links := link(
		[2]string{"start", "t"}, [2]string{"n", "e"}, [2]string{"a", "m"},
		[2]string{"A", "c"}, [2]string{"0", "o"}, [2]string{"E", "a"},
		[2]string{"k", "end"}, [2]string{"start", "h"}, [2]string{"o", "n"},
		[2]string{"m", "end"}, [2]string{"t", "E"}, [2]string{"start", "0"},
		[2]string{"h", "A"}, [2]string{"e", "end"}, [2]string{"c", "k"},
		[2]string{"n", "m"}, [2]string{"h", "n"},
	)
	a := Solve(links, "start", "end", 10)
	if a == nil {
		t.Fatal("expected an assignment, got nil")
	}
	if got := turnsUsed(a); got != 8 {
		t.Errorf("turnsUsed = %d, want 8 (3 disjoint 5-hop paths, 10 ants split 4/3/3)", got)
	}
}

func TestSolve_AvoidsUnnecessaryDetourPaths(t *testing.T) {
	links := link(
		[2]string{"start", "short"}, [2]string{"short", "end"},

		[2]string{"start", "mid"}, [2]string{"mid", "shortcut"},
		[2]string{"shortcut", "end"},
		[2]string{"mid", "detourA"}, [2]string{"detourA", "detourB"}, [2]string{"detourB", "end"},

		[2]string{"start", "long1"}, [2]string{"long1", "long2"}, [2]string{"long2", "bridge"},
		[2]string{"bridge", "shortcut"}, 
	)

	a := Solve(links, "start", "end", 3)
	if a == nil {
		t.Fatal("expected an assignment, got nil")
	}

	usesLong := false
	for _, p := range a.Paths {
		for _, room := range p {
			if room == "long1" {
				usesLong = true
			}
		}
	}
	if usesLong {
		t.Errorf("solver used the long/bridge branch even though it wasn't needed to hit the turn floor: %v", a.Paths)
	}
}

func TestSolve_EveryAntAssigned(t *testing.T) {
	links := link([2]string{"start", "a"}, [2]string{"a", "b"}, [2]string{"b", "end"})
	a := Solve(links, "start", "end", 7)
	if a == nil {
		t.Fatal("expected an assignment, got nil")
	}
	seen := make([]bool, len(a.Paths))
	for _, p := range a.AntToPath {
		if p < 0 || p >= len(a.Paths) {
			t.Fatalf("AntToPath contains out-of-range path index %d", p)
		}
		seen[p] = true
	}
	if len(a.AntToPath) != 7 {
		t.Errorf("len(AntToPath) = %d, want 7 (one entry per ant)", len(a.AntToPath))
	}
}
