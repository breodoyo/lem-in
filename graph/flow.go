package graph

import "sort"

type flowEdge struct {
	to   int
	cap  int
	orig int
}

type flowGraph struct {
	adj   [][]int
	edges []flowEdge
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{adj: make([][]int, n)}
}

func (g *flowGraph) addEdge(u, v, cap int) {
	g.edges = append(g.edges, flowEdge{to: v, cap: cap, orig: cap})
	g.edges = append(g.edges, flowEdge{to: u, cap: 0, orig: 0})
	g.adj[u] = append(g.adj[u], len(g.edges)-2)
	g.adj[v] = append(g.adj[v], len(g.edges)-1)
}

func (g *flowGraph) maxFlow(s, t int, limit int) {
	pushed := 0
	for limit < 0 || pushed < limit {
		parentEdge := make([]int, len(g.adj))
		for i := range parentEdge {
			parentEdge[i] = -1
		}
		parentEdge[s] = -2
		queue := []int{s}

		for len(queue) > 0 && parentEdge[t] == -1 {
			u := queue[0]
			queue = queue[1:]
			for _, eIdx := range g.adj[u] {
				e := g.edges[eIdx]
				if parentEdge[e.to] == -1 && e.cap > 0 {
					parentEdge[e.to] = eIdx
					queue = append(queue, e.to)
				}
			}
		}

		if parentEdge[t] == -1 {
			return
		}

		bottleneck := int(^uint(0) >> 1)
		for v := t; v != s; {
			eIdx := parentEdge[v]
			if g.edges[eIdx].cap < bottleneck {
				bottleneck = g.edges[eIdx].cap
			}
			v = g.edges[eIdx^1].to
		}
		for v := t; v != s; {
			eIdx := parentEdge[v]
			g.edges[eIdx].cap -= bottleneck
			g.edges[eIdx^1].cap += bottleneck
			v = g.edges[eIdx^1].to
		}
		pushed += bottleneck
	}
}

func (g *flowGraph) decomposePaths(s, t int, roomOf map[int]string) [][]string {
	used := make([][]int, len(g.adj))
	for u := range g.adj {
		for _, eIdx := range g.adj[u] {
			e := g.edges[eIdx]
			if e.orig > 0 && e.cap < e.orig {
				used[u] = append(used[u], eIdx)
			}
		}
	}

	var paths [][]string
	for len(used[s]) > 0 {
		var nodePath []int
		visited := map[int]bool{}
		cur := s
		nodePath = append(nodePath, cur)
		visited[cur] = true

		for cur != t && len(used[cur]) > 0 {
			eIdx := used[cur][len(used[cur])-1]
			used[cur] = used[cur][:len(used[cur])-1]
			cur = g.edges[eIdx].to
			if visited[cur] {
				break 
			}
			nodePath = append(nodePath, cur)
			visited[cur] = true
		}
		if cur != t {
			continue 
		}

		var roomPath []string
		for _, nd := range nodePath {
			name := roomOf[nd]
			if len(roomPath) == 0 || roomPath[len(roomPath)-1] != name {
				roomPath = append(roomPath, name) 
			}
		}
		paths = append(paths, roomPath)
	}
	return paths
}

func VertexDisjointPaths(links map[string][]string, start, end string) [][]string {
	return VertexDisjointPathsLimited(links, start, end, -1)
}

func VertexDisjointPathsLimited(links map[string][]string, start, end string, limit int) [][]string {
	rooms := map[string]bool{start: true, end: true}
	for room, neighbors := range links {
		rooms[room] = true
		for _, n := range neighbors {
			rooms[n] = true
		}
	}

	ids := map[string][2]int{}
	next := 0
	nodeIDs := func(room string) (int, int) {
		if v, ok := ids[room]; ok {
			return v[0], v[1]
		}
		var in, out int
		if room == start || room == end {
			in, out = next, next
			next++
		} else {
			in, out = next, next+1
			next += 2
		}
		ids[room] = [2]int{in, out}
		return in, out
	}

	sortedRooms := make([]string, 0, len(rooms))
	for room := range rooms {
		sortedRooms = append(sortedRooms, room)
	}
	sort.Strings(sortedRooms)
	for _, room := range sortedRooms {
		nodeIDs(room)
	}

	g := newFlowGraph(next)

	for _, room := range sortedRooms {
		if room == start || room == end {
			continue
		}
		in, out := nodeIDs(room)
		g.addEdge(in, out, 1)
	}

	sortedLinkRooms := make([]string, 0, len(links))
	for room := range links {
		sortedLinkRooms = append(sortedLinkRooms, room)
	}
	sort.Strings(sortedLinkRooms)

	seen := map[[2]string]bool{}
	for _, a := range sortedLinkRooms {
		neighbors := append([]string(nil), links[a]...)
		sort.Strings(neighbors)
		for _, b := range neighbors {
			if seen[[2]string{a, b}] {
				continue
			}
			seen[[2]string{a, b}] = true
			seen[[2]string{b, a}] = true

			_, aOut := nodeIDs(a)
			bIn, _ := nodeIDs(b)
			g.addEdge(aOut, bIn, 1)

			_, bOut := nodeIDs(b)
			aIn, _ := nodeIDs(a)
			g.addEdge(bOut, aIn, 1)
		}
	}

	_, source := nodeIDs(start)
	sink, _ := nodeIDs(end)
	g.maxFlow(source, sink, limit)

	roomOf := make(map[int]string, len(ids)*2)
	for room, v := range ids {
		roomOf[v[0]] = room
		roomOf[v[1]] = room
	}

	return g.decomposePaths(source, sink, roomOf)
}