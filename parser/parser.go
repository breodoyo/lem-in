package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrInvalidFormat = errors.New("ERROR: invalid data format")
//takes a files destination and return a successfully fiiled colony structure or an error
func ParseFile(path string) (*Colony, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w, could not read file: %v", ErrInvalidFormat, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w, could not read file", ErrInvalidFormat)
	}

	return Parse(lines)
}

func Parse(lines []string) (*Colony, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("%w, empty file", ErrInvalidFormat)
	}

	colony := NewColony()
	colony.RawLines = lines

	i := 0
	numAntsSet := false
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##")) {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("%w, invalid number of ants", ErrInvalidFormat)
		}
		colony.NumAnts = n
		numAntsSet = true
		i++
		break
	}
	if !numAntsSet {
		return nil, fmt.Errorf("%w, invalid number of ants", ErrInvalidFormat)
	}

	nextIsStart := false
	nextIsEnd := false

	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		if line == "##start" {
			if nextIsStart || nextIsEnd {
				return nil, fmt.Errorf("%w, malformed start/end marker", ErrInvalidFormat)
			}
			nextIsStart = true
			continue
		}
		if line == "##end" {
			if nextIsStart || nextIsEnd {
				return nil, fmt.Errorf("%w, malformed start/end marker", ErrInvalidFormat)
			}
			nextIsEnd = true
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return nil, fmt.Errorf("%w, malformed link %q", ErrInvalidFormat, line)
			}
			a, b := parts[0], parts[1]
			if err := validateLink(colony, a, b); err != nil {
				return nil, err
			}
			colony.AddLink(a, b)
			continue
		}

		room, err := parseRoom(line)
		if err != nil {
			return nil, err
		}
		if _, exists := colony.Rooms[room.Name]; exists {
			return nil, fmt.Errorf("%w, duplicate room %q", ErrInvalidFormat, room.Name)
		}
		colony.Rooms[room.Name] = room

		if nextIsStart {
			if colony.Start != "" {
				return nil, fmt.Errorf("%w, multiple start rooms", ErrInvalidFormat)
			}
			colony.Start = room.Name
			nextIsStart = false
		}
		if nextIsEnd {
			if colony.End != "" {
				return nil, fmt.Errorf("%w, multiple end rooms", ErrInvalidFormat)
			}
			colony.End = room.Name
			nextIsEnd = false
		}
	}

	if colony.Start == "" {
		return nil, fmt.Errorf("%w, no start room found", ErrInvalidFormat)
	}
	if colony.End == "" {
		return nil, fmt.Errorf("%w, no end room found", ErrInvalidFormat)
	}
	if colony.Start == colony.End {
		return nil, fmt.Errorf("%w, start and end room are the same", ErrInvalidFormat)
	}

	return colony, nil
}

func parseRoom(line string) (*Room, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return nil, fmt.Errorf("%w, malformed room %q", ErrInvalidFormat, line)
	}

	name := fields[0]
	if name == "" || strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		return nil, fmt.Errorf("%w, invalid room name %q", ErrInvalidFormat, name)
	}

	x, errX := strconv.Atoi(fields[1])
	y, errY := strconv.Atoi(fields[2])
	if errX != nil || errY != nil {
		return nil, fmt.Errorf("%w, invalid coordinates for room %q", ErrInvalidFormat, name)
	}

	return &Room{Name: name, X: x, Y: y}, nil
}

func validateLink(c *Colony, a, b string) error {
	if a == b {
		return fmt.Errorf("%w, room %q links to itself", ErrInvalidFormat, a)
	}
	if _, ok := c.Rooms[a]; !ok {
		return fmt.Errorf("%w, link references unknown room %q", ErrInvalidFormat, a)
	}
	if _, ok := c.Rooms[b]; !ok {
		return fmt.Errorf("%w, link references unknown room %q", ErrInvalidFormat, b)
	}
	if c.HasLink(a, b) {
		return fmt.Errorf("%w, duplicate link %q-%q", ErrInvalidFormat, a, b)
	}
	return nil
}
