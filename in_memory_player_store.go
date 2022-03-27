package poker

type InMemoryPlayerScore struct {
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerScore {
	return &InMemoryPlayerScore{map[string]int{}}
}

func (i *InMemoryPlayerScore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerScore) GetLeague() (league League) {
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}

	return
}

func (i *InMemoryPlayerScore) RecordWin(name string) {
	i.store[name]++
}
