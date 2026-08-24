# lem-in

A digital ant farm simulation and network flow optimization engine written in Go. The program finds the quickest combinations of paths to guide $N$ ants across an interconnected network of rooms and tunnels from a designated starting point (`##start`) to an exit (`##end`) in the fewest number of turns.

---

## Overview

The goal of `lem-in` is to navigate a population of ants through an anthill graph with minimal traffic congestion. Every room along the path can only hold **one ant at a time** (excluding the start and end rooms), meaning the absolute shortest path is not always the most efficient choice when multiple ants are waiting. The engine must dynamically calculate and distribute ants across multiple non-overlapping paths to avoid bottlenecks.

This project is a deep dive into graph theory, network flow algorithms (such as Edmonds-Karp, Suurballe's, or Dinic's), structural data parsing, and stream orchestration.

---

## Project Structure

The project follows clean Go design principles and is structured as a modular codebase:

```text
.
├── go.mod            # Go module definitions
├── main.go           # Application entrypoint
├── main_test.go      # Integration and regression tests
├── README.md         # Documentation
├── graph/            # Adjacency structures, path containers, and residual graph math
│   ├── flow.go
│   ├── flow_test.go
│   └── graph.go
├── parser/           # Validates constraints, handles strings, tokens, and file input
│   ├── parser.go
│   ├── parser_test.go
│   └── types.go
├── simulator/        # Manages turn-by-turn ant dispatches and path-walking routines
│   └── simulator.go
├── solver/           # Multi-path optimization algorithms (disjoint path selection)
│   ├── solver.go
│   └── solver_test.go
├── testdata/         # Collection of valid map samples and invalid edge-case matrices
│   ├── badexample00.txt
│   └── example01.txt
└── verify/           # Validation layers confirming output compliance with runtime rules
    └── verify.go
```

---

## Core Rules & Constraints

### The Environment
- **Rooms**: Defined by a unique name and coordinates (`name coord_x coord_y`). 
  - Room names will never start with `L` or `#` and contain no spaces.
  - Coordinates are always integers (`int`).
- **Tunnels**: Link exactly two rooms together (`room1-room2`). 
  - Two rooms can have at most one tunnel connecting them.
  - A room can connect to multiple other rooms.

### The Logic
- At the start, all ants wait in the `##start` room. The target is the `##end` room.
- Each room can hold **only one ant at any given time** (except `##start` and `##end`).
- Each tunnel can handle **only one ant per turn**.
- During a turn, you can only move an ant if its destination room is empty or will be vacated during the same turn.
- Ants must be routed to minimize total turns.

### Error Handling
The program must handle errors gracefully without crashing unexpectedly. If the input contains loops, bad coordinates, duplicate rooms, missing links, no exit, or an invalid number of ants, the program prints a formatted message to standard error:
`ERROR: invalid data format` *(or specific tracking like `ERROR: invalid data format, no start room found`)*.

---

## Input & Output Formats

### Input Map File Format
```text
number_of_ants
the_rooms
the_links
```

### Output Protocol
The program echoes the valid map format, followed by an empty line, followed by the turn-by-turn simulation sequences:
```text
L[ant_number]-[room_name] L[ant_number]-[room_name] ...
```

---

## Usage

Only standard Go packages are used in this engine. Run the builder or execution pipe directly:

```bash
go run . <path_to_map_file>
```

### Example 1
```bash
$ go run . test0.txt
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5

L1-3 L2-2
L1-4 L2-5 L3-3
L1-0 L2-6 L3-4
L2-0 L3-0
```

### Example 2
```bash
$ go run . test1.txt
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

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3
L3-1
```

---

## Testing

Unit tests ensure parser limits and node routing equations operate properly. Execute your test suites via:

```bash
# Run all project tests
go test ./... -v
```

---

## Simulation Visualizer

The program outputs linear stream data which can easily pipe into a graphical interpreter. The spatial properties (`coord_x` and `coord_y`) are designed specifically to render visual structures here.

```bash
./lem-in ant-farm.txt | ./visualizer
```