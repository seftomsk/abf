package storage

const (
	white = "white"
	black = "black"
)

type InMemory struct {
	collection map[string]map[string][]string
}

func (s *InMemory) AddToWhiteList(ip, mask string) {
	s.collection[white][mask] = append(s.collection[white][mask], ip)
}

func (s *InMemory) AddToBlackList(ip, mask string) {
	s.collection[black][mask] = append(s.collection[black][mask], ip)
}

func (s *InMemory) GetAll() map[string]map[string][]string {
	return s.collection
}

func NewInMemory() *InMemory {
	collection := make(map[string]map[string][]string)
	collection[white] = make(map[string][]string)
	collection[black] = make(map[string][]string)

	return &InMemory{collection: collection}
}
