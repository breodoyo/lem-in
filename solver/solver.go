package solver

import (
	"sort"

	"lem-in/graph"
)

type Assignment struct {
	Paths     [][]string 
	AntToPath []int      
}

func Solve(links map[string][]string, start, end string, numAnts int) *Assignment {
	maxPaths := len(graph.VertexDisjointPaths(links, start, end))
	if maxPaths == 0 {
		return nil
	}

	var candidates []*Assignment
	var turnsFor []int

	for k := 1; k <= maxPaths; k++ {
		paths := graph.VertexDisjointPathsLimited(links, start, end, k)
		sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })

		a, turns := assign(paths, numAnts)
		candidates = append(candidates, a)
		turnsFor = append(turnsFor, turns)
	}

	bestTurns := turnsFor[0]
	for _, t := range turnsFor {
		if t < bestTurns {
			bestTurns = t
		}
	}
	for i, t := range turnsFor {
		if t == bestTurns {
			return candidates[i]
		}
	}
	return candidates[len(candidates)-1]
}

func assign(paths [][]string, numAnts int) (*Assignment, int) {
	counts := make([]int, len(paths))
	antToPath := make([]int, numAnts)

	hops := make([]int, len(paths))
	for p, path := range paths {
		hops[p] = len(path) - 1
	}

	turns := 0
	for ant := 0; ant < numAnts; ant++ {
		best := 0
		bestFinish := hops[0] + counts[0]
		for p := 1; p < len(paths); p++ {
			finish := hops[p] + counts[p]
			if finish < bestFinish {
				bestFinish = finish
				best = p
			}
		}
		counts[best]++
		antToPath[ant] = best
		if bestFinish > turns {
			turns = bestFinish
		}
	}

	return &Assignment{Paths: paths, AntToPath: antToPath}, turns
}