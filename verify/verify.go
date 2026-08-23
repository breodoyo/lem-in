package verify

import (
	"fmt"
	"strconv"
	"strings"

	"lem-in/parser"
)

func Check(c *parser.Colony, numAnts int, turns []string) error {
	pos := make([]string, numAnts+1)
	for i := 1; i <= numAnts; i++ {
		pos[i] = c.Start
	}
	occupant := map[string]int{} 

	for turnIdx, line := range turns {
		type move struct {
			ant  int
			cur  string
			dest string
		}
		var moves []move
		movedThisTurn := map[int]bool{}

		for _, tok := range strings.Fields(line) {
			if !strings.HasPrefix(tok, "L") {
				return fmt.Errorf("turn %d: malformed move %q", turnIdx+1, tok)
			}
			rest := tok[1:]
			dash := strings.IndexByte(rest, '-')
			if dash < 0 {
				return fmt.Errorf("turn %d: malformed move %q", turnIdx+1, tok)
			}
			antNum, err := strconv.Atoi(rest[:dash])
			if err != nil || antNum < 1 || antNum > numAnts {
				return fmt.Errorf("turn %d: invalid ant number in %q", turnIdx+1, tok)
			}
			dest := rest[dash+1:]

			if movedThisTurn[antNum] {
				return fmt.Errorf("turn %d: ant %d moved twice in one turn", turnIdx+1, antNum)
			}
			movedThisTurn[antNum] = true

			cur := pos[antNum]
			if cur == c.End {
				return fmt.Errorf("turn %d: ant %d moved after already reaching end", turnIdx+1, antNum)
			}
			if !c.HasLink(cur, dest) {
				return fmt.Errorf("turn %d: ant %d has no tunnel from %q to %q", turnIdx+1, antNum, cur, dest)
			}

			moves = append(moves, move{ant: antNum, cur: cur, dest: dest})
		}

		next := make(map[string]int, len(occupant))
		for room, ant := range occupant {
			next[room] = ant
		}
		for _, m := range moves {
			if m.cur != c.Start && m.cur != c.End {
				delete(next, m.cur)
			}
		}
		for _, m := range moves {
			if m.dest == c.Start || m.dest == c.End {
				continue 
			}
			if other, taken := next[m.dest]; taken && other != m.ant {
				return fmt.Errorf("turn %d: room %q double-booked by ants %d and %d", turnIdx+1, m.dest, other, m.ant)
			}
			next[m.dest] = m.ant
		}
		occupant = next

		for _, m := range moves {
			pos[m.ant] = m.dest
		}
	}

	for i := 1; i <= numAnts; i++ {
		if pos[i] != c.End {
			return fmt.Errorf("ant %d never reached end (last seen at %q)", i, pos[i])
		}
	}
	return nil
}
