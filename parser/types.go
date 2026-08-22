package parser

type Room struct {
	Name string
	X, Y int
}

type Colony struct {
	NumAnts int
	Rooms   map[string]*Room
	Links   map[string][]string
	Start   string
	End     string

	RawLines []string
}

func NewColony() *Colony {
	return &Colony{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
	}
}

func (c *Colony) AddLink(a, b string) {
	c.Links[a] = append(c.Links[a], b)
	c.Links[b] = append(c.Links[b], a)
}

func (c *Colony) HasLink(a, b string) bool {
	for _, n := range c.Links[a] {
		if n == b {
			return true
		}
	}
	return false
}
