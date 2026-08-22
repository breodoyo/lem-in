package solver

import (
	"sort"

	"lem-in/graph"
)

type Assignment struct {
	Paths       [][]string
	AntToPath   []int      
}

func Solve(links map[string][]string, start, end string, numAnts int) *Assignment {
	g := graph.New(links)
	candidates := g.AllEdgeDisjointPaths(start, end)
	if len(candidates) == 0 {
		return nil
	}

	paths := roomDisjointSubset(candidates, start, end)

	
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

func roomDisjointSubset(candidates [][]string, start, end string) [][]string {
	sort.Slice(candidates, func(i, j int) bool { return len(candidates[i]) < len(candidates[j]) })

	usedRooms := map[string]bool{}
	var chosen [][]string

	for _, path := range candidates {
		if roomsFree(path, start, end, usedRooms) {
			chosen = append(chosen, path)
			for _, room := range path {
				if room != start && room != end {
					usedRooms[room] = true
				}
			}
		}
	}

	return chosen
}

func roomsFree(path []string, start, end string, used map[string]bool) bool {
	for _, room := range path {
		if room == start || room == end {
			continue
		}
		if used[room] {
			return false
		}
	}
	return true
}
