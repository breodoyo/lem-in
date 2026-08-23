package graph

import "testing"

func link(pairs ...[2]string) map[string][]string {
	m := map[string][]string{}
	for _, p := range pairs {
		m[p[0]] = append(m[p[0]], p[1])
		m[p[1]] = append(m[p[1]], p[0])
	}
	return m
}

func TestVertexDisjointPaths_SinglePath(t *testing.T) {
	links := link([2]string{"start", "a"}, [2]string{"a", "end"})
	paths := VertexDisjointPaths(links, "start", "end")
	if len(paths) != 1 {
		t.Fatalf("got %d paths, want 1: %v", len(paths), paths)
	}
	if len(paths[0]) != 3 {
		t.Errorf("path length = %d, want 3 (start,a,end)", len(paths[0]))
	}
}

func TestVertexDisjointPaths_NoPath(t *testing.T) {
	links := link([2]string{"start", "a"}, [2]string{"b", "end"})
	paths := VertexDisjointPaths(links, "start", "end")
	if len(paths) != 0 {
		t.Fatalf("got %d paths, want 0 (start and end are disconnected): %v", len(paths), paths)
	}
}

func TestVertexDisjointPaths_TwoDisjoint(t *testing.T) {
	// start-1-2-end and a direct start-3-end.
	links := link(
		[2]string{"start", "1"}, [2]string{"1", "2"}, [2]string{"2", "end"},
		[2]string{"start", "3"}, [2]string{"3", "end"},
	)
	paths := VertexDisjointPaths(links, "start", "end")
	if len(paths) != 2 {
		t.Fatalf("got %d paths, want 2: %v", len(paths), paths)
	}
}
func TestVertexDisjointPaths_BridgeMonopoly(t *testing.T) {
	links := link(
		[2]string{"start", "t"}, [2]string{"n", "e"}, [2]string{"a", "m"},
		[2]string{"A", "c"}, [2]string{"0", "o"}, [2]string{"E", "a"},
		[2]string{"k", "end"}, [2]string{"start", "h"}, [2]string{"o", "n"},
		[2]string{"m", "end"}, [2]string{"t", "E"}, [2]string{"start", "0"},
		[2]string{"h", "A"}, [2]string{"e", "end"}, [2]string{"c", "k"},
		[2]string{"n", "m"}, [2]string{"h", "n"},
	)
	paths := VertexDisjointPaths(links, "start", "end")
	if len(paths) != 3 {
		t.Fatalf("got %d paths, want 3 (the graph's true max-flow): %v", len(paths), paths)
	}
	for _, p := range paths {
		if len(p) != 6 { // 5 hops: start + 4 intermediate + end
			t.Errorf("path %v has length %d, want 6", p, len(p))
		}
	}
}

func TestVertexDisjointPaths_Deterministic(t *testing.T) {
	links := link(
		[2]string{"0", "1"}, [2]string{"1", "4"}, [2]string{"4", "5"},
		[2]string{"2", "4"}, [2]string{"0", "2"},
		[2]string{"3", "0"}, [2]string{"4", "3"},
	)
	first := VertexDisjointPaths(links, "0", "5")
	for i := 0; i < 20; i++ {
		got := VertexDisjointPaths(links, "0", "5")
		if len(got) != len(first) || (len(got) > 0 && !equalPaths(got[0], first[0])) {
			t.Fatalf("run %d produced a different result than the first run: %v vs %v", i, got, first)
		}
	}
}

func equalPaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
