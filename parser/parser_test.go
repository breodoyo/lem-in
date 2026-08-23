package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParse_Valid(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(`
3
##start
0 1 0
##end
1 5 0
2 9 0
3 13 0
0-2
2-3
3-1
`), "\n")

	c, err := Parse(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.NumAnts != 3 {
		t.Errorf("NumAnts = %d, want 3", c.NumAnts)
	}
	if c.Start != "0" {
		t.Errorf("Start = %q, want %q", c.Start, "0")
	}
	if c.End != "1" {
		t.Errorf("End = %q, want %q", c.End, "1")
	}
	if len(c.Rooms) != 4 {
		t.Errorf("len(Rooms) = %d, want 4", len(c.Rooms))
	}
	if !c.HasLink("0", "2") || !c.HasLink("2", "3") || !c.HasLink("3", "1") {
		t.Errorf("expected links 0-2, 2-3, 3-1 to all be present")
	}
}

func TestParse_StripsComments(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(`
1
#a comment
##start
0 0 0
##end
1 1 1
0-1
`), "\n")

	c, err := Parse(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Rooms) != 2 {
		t.Errorf("len(Rooms) = %d, want 2 (comment line should not become a room)", len(c.Rooms))
	}
}
func TestParse_HashLinesAreAlwaysComments(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(`
1
##start
0 0 0
##end
1 1 1
0-1
#1 1 1
`), "\n")

	c, err := Parse(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Rooms) != 2 {
		t.Errorf("len(Rooms) = %d, want 2 (the trailing #-line should be dropped as a comment, not parsed as a room)", len(c.Rooms))
	}
}

func TestParse_InvalidCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty file", ``},
		{"missing ant count", `##start
0 0 0
##end
1 1 1
0-1`},
		{"zero ants", `0
##start
0 0 0
##end
1 1 1
0-1`},
		{"negative ants", `-3
##start
0 0 0
##end
1 1 1
0-1`},
		{"no start room", `1
0 0 0
##end
1 1 1
0-1`},
		{"no end room", `1
##start
0 0 0
1 1 1
0-1`},
		{"duplicate start", `1
##start
0 0 0
##start
1 1 1
##end
2 2 2
0-1
1-2`},
		{"start marker immediately followed by end marker", `1
##start
##end
0 0 0
1 1 1
0-1`},
		{"duplicate room name", `1
##start
0 0 0
##end
1 1 1
0 2 2
0-1`},
		{"link to unknown room", `1
##start
0 0 0
##end
1 1 1
0-9`},
		{"self link", `1
##start
0 0 0
##end
1 1 1
0-0`},
		{"duplicate link", `1
##start
0 0 0
##end
1 1 1
0-1
0-1`},
		{"room name starts with L", `1
##start
L0 0 0
##end
1 1 1
L0-1`},
		{"non-integer coordinates", `1
##start
0 a b
##end
1 1 1
0-1`},
		{"malformed room line (wrong field count)", `1
##start
0 0
##end
1 1 1
0-1`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var lines []string
			if tc.input != "" {
				lines = strings.Split(tc.input, "\n")
			}
			_, err := Parse(lines)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !errors.Is(err, ErrInvalidFormat) {
				t.Errorf("error %q does not wrap ErrInvalidFormat", err)
			}
		})
	}
}
