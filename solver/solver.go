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
	paths := graph.VertexDisjointPaths(links, start, end)
	if len(paths) == 0 {
		return nil
	}

	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })

	counts := make([]int, len(paths))
	antToPath := make([]int, numAnts)

	hops := make([]int, len(paths))
	for p, path := range paths {
		hops[p] = len(path) - 1
	}

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
	}

	return &Assignment{Paths: paths, AntToPath: antToPath}
}