package main

import (
	"errors"
	"math/rand"
	"time"
)

const (
	DESK_SIZE  int = 144
	MAGIC_BASE int = 7
)

type Deck struct {
	Size   int
	header int
	tail   int
	Magic  int
	Stack  []string
	rnd *rand.Rand
}

func (s *Deck) New() {
	if s.Magic <= 0 {
		s.Magic = MAGIC_BASE
	}
	if s.Size <= 0 {
		s.Size = DESK_SIZE
	}
	s.Stack = make([]string, 0)
	for i := range 4 {
		s.Stack = append(s.Stack, BAMBOO1, BAMBOO2, BAMBOO3, BAMBOO4, BAMBOO5, BAMBOO6, BAMBOO7, BAMBOO8, BAMBOO9)
		s.Stack = append(s.Stack, CHARACTER1, CHARACTER2, CHARACTER3, CHARACTER4, CHARACTER5, CHARACTER6, CHARACTER7, CHARACTER8, CHARACTER9)
		s.Stack = append(s.Stack, DOTS1, DOTS2, DOTS3, DOTS4, DOTS5, DOTS6, DOTS7, DOTS8, DOTS9)
		s.Stack = append(s.Stack, EAST, SOUTH, WEST, NORTH, RED, GREEN, WHITE)
		if i == 0 {
			s.Stack = append(s.Stack, F_BAMBOO, F_CHRYSANTHEMUM, F_ORCHID, F_PLUMBLOSSOM, F_SPRING, F_SUMMER, F_AUTUMN, F_WINTER)
		}
	}
	src := rand.NewSource(time.Now().UnixNano())
	s.rnd = rand.New(src)
	
}
func (s *Deck) Shuffle() {
	s.header = 0
	cut := s.rnd.Intn(s.Magic) + 1
	s.tail = s.Size - cut
	s.rnd.Shuffle(s.Size, func(i, j int) {
		s.Stack[i], s.Stack[j] = s.Stack[j], s.Stack[i]
	})
}

func (s *Deck) Draw() (Tile, error) {
	t := Tile{}
	if s.header > s.tail {
		return t, errors.New("no more draw")
	}
	t.From(s.Stack[s.header])
	s.header++
	return t, nil
}

func (s *Deck) Kong() (Tile, error) {
	t := Tile{}
	if s.tail < s.header {
		return t, errors.New("no more knog")
	}
	t.From(s.Stack[s.tail])
	s.tail--
	return t, nil
}

func (s *Deck) Dice() []int {
	
	dt := []int{0, 0}
	dt[0] = s.rnd.Intn(6) + 1
	dt[1] = s.rnd.Intn(6) + 1
	return dt
}
