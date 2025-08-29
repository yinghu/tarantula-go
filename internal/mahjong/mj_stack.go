package mahjong

import (
	"errors"
	"math/rand"
	"time"
)

const (
	DESK_SIZE  int = 144
	MAGIC_BASE int = 7
)

type Stack struct {
	Size   int
	header int
	tail   int
	Magic  int
	Deck   []string
}

func (s *Stack) New() {
	if s.Magic <= 0 {
		s.Magic = MAGIC_BASE
	}
	if s.Size <= 0 {
		s.Size = DESK_SIZE
	}
	s.Deck = make([]string, 0)
	for i := range 4 {
		s.Deck = append(s.Deck, BAMBOO1, BAMBOO2, BAMBOO3, BAMBOO4, BAMBOO5, BAMBOO6, BAMBOO7, BAMBOO8, BAMBOO9)
		s.Deck = append(s.Deck, CHARACTER1, CHARACTER2, CHARACTER3, CHARACTER4, CHARACTER5, CHARACTER6, CHARACTER7, CHARACTER8, CHARACTER9)
		s.Deck = append(s.Deck, DOTS1, DOTS2, DOTS3, DOTS4, DOTS5, DOTS6, DOTS7, DOTS8, DOTS9)
		s.Deck = append(s.Deck, EAST, SOUTH, WEST, NORTH, RED, GREEN, WHITE)
		if i == 0 {
			s.Deck = append(s.Deck, F_BAMBOO, F_CHRYSANTHEMUM, F_ORCHID, F_PLUMBLOSSOM, F_SPRING, F_SUMMER, F_AUTUMN, F_WINTER)
		}
	}
}
func (s *Stack) Shuffle() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	s.header = 0
	cut := rnd.Intn(s.Magic) + 1
	s.tail = s.Size - cut
	rnd.Shuffle(s.Size, func(i, j int) {
		s.Deck[i], s.Deck[j] = s.Deck[j], s.Deck[i]
	})
}

func (s *Stack) Draw() (Tile, error) {
	t := Tile{}
	if s.header > s.tail {
		return t, errors.New("no more draw")
	}
	t.From(s.Deck[s.header])
	s.header++
	return t, nil
}

func (s *Stack) Kong() (Tile, error) {
	t := Tile{}
	if s.tail < s.header {
		return t, errors.New("no more knog")
	}
	t.From(s.Deck[s.tail])
	s.tail--
	return t, nil
}
