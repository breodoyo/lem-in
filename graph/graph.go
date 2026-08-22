package graph

type Graph struct {
	adj map[string]map[string]bool
}

func New(links map[string][]string) *Graph {
	g := &Graph{adj: make(map[string]map[string]bool)}
	for room, neighbors := range links {
		if g.adj[room] == nil {
			g.adj[room] = make(map[string]bool)
		}
		for _, n := range neighbors {
			g.adj[room][n] = true
		}
	}
	return g
}

func (g *Graph) RemoveEdge(a, b string) {
	delete(g.adj[a], b)
	delete(g.adj[b], a)
}

func (g *Graph) BFSShortestPath(start, end string) []string {
	if start == end {
		return []string{start}
	}

	visited := map[string]bool{start: true}
	prev := map[string]string{}
	queue := []string{start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == end {
			return reconstructPath(prev, start, end)
		}

		neighbors := make([]string, 0, len(g.adj[cur]))
		for n := range g.adj[cur] {
			neighbors = append(neighbors, n)
		}

		for _, n := range neighbors {
			if !visited[n] {
				visited[n] = true
				prev[n] = cur
				queue = append(queue, n)
			}
		}
	}

	return nil
}

func reconstructPath(prev map[string]string, start, end string) []string {
	path := []string{end}
	cur := end
	for cur != start {
		cur = prev[cur]
		path = append([]string{cur}, path...)
	}
	return path
}

func (g *Graph) AllEdgeDisjointPaths(start, end string) [][]string {
	working := g.clone()
	var paths [][]string

	for {
		path := working.BFSShortestPath(start, end)
		if path == nil {
			break
		}
		paths = append(paths, path)
		for i := 0; i < len(path)-1; i++ {
			working.RemoveEdge(path[i], path[i+1])
		}
	}

	return paths
}

func (g *Graph) clone() *Graph {
	cp := &Graph{adj: make(map[string]map[string]bool, len(g.adj))}
	for room, neighbors := range g.adj {
		cp.adj[room] = make(map[string]bool, len(neighbors))
		for n := range neighbors {
			cp.adj[room][n] = true
		}
	}
	return cp
}
