package simulator

import (
	"fmt"
	"strings"

	"lem-in/solver"
)

//tracks one ant's progress along its assigned path.
type antState struct {
	pathIdx int 
	pos     int
}

func Run(a *solver.Assignment, numAnts int) []string {
	ants := make([]*antState, numAnts)
	for i := 0; i < numAnts; i++ {
		ants[i] = &antState{pathIdx: a.AntToPath[i], pos: 0}
	}

	queues := make([][]int, len(a.Paths))
	for i := 0; i < numAnts; i++ {
		p := ants[i].pathIdx
		queues[p] = append(queues[p], i)
	}

	entered := make([]int, len(a.Paths))

	var lines []string
	finished := 0

	for finished < numAnts {
		type move struct {
			ant int
			str string
		}
		var turnMoves []move

		for p, path := range a.Paths {
			hops := len(path) - 1

			
			for idx := entered[p] - 1; idx >= 0; idx-- {
				antNum := queues[p][idx]
				st := ants[antNum]
				if st.pos >= hops {
					continue
				}
				nextPos := st.pos + 1
				roomAhead := path[nextPos]
				
				aheadIsFree := nextPos == hops
				if !aheadIsFree {
					aheadIsFree = true
					for j := idx + 1; j < entered[p]; j++ {
						other := ants[queues[p][j]]
						if other.pos == nextPos {
							aheadIsFree = false
							break
						}
					}
				}
				if aheadIsFree {
					st.pos = nextPos
					turnMoves = append(turnMoves, move{ant: antNum + 1, str: fmt.Sprintf("L%d-%s", antNum+1, roomAhead)})
					if st.pos == hops {
						finished++
					}
				}
			}

			if entered[p] < len(queues[p]) {
				firstRoomFree := true
				for j := 0; j < entered[p]; j++ {
					if ants[queues[p][j]].pos == 1 {
						firstRoomFree = false
						break
					}
				}
				if firstRoomFree || hops == 1 {
					antNum := queues[p][entered[p]]
					st := ants[antNum]
					st.pos = 1
					turnMoves = append(turnMoves, move{ant: antNum + 1, str: fmt.Sprintf("L%d-%s", antNum+1, path[1])})
					entered[p]++
					if st.pos == hops {
						finished++
					}
				}
			}
		}

		if len(turnMoves) == 0 {
			break
		}

		for i := 0; i < len(turnMoves); i++ {
			for j := i + 1; j < len(turnMoves); j++ {
				if turnMoves[j].ant < turnMoves[i].ant {
					turnMoves[i], turnMoves[j] = turnMoves[j], turnMoves[i]
				}
			}
		}
		strs := make([]string, len(turnMoves))
		for i, m := range turnMoves {
			strs[i] = m.str
		}
		lines = append(lines, strings.Join(strs, " "))
	}

	return lines
}
