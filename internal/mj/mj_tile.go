package mj

import (
	"strconv"
)

const (
	FLOWER    string = "F"
	BAMBOO    string = "B"
	CHARACTER string = "C"
	DOTS      string = "D"
	HORNOR    string = "H"

	BAMBOO1 string = "B1"
	BAMBOO2 string = "B2"
	BAMBOO3 string = "B3"
	BAMBOO4 string = "B4"
	BAMBOO5 string = "B5"
	BAMBOO6 string = "B6"
	BAMBOO7 string = "B7"
	BAMBOO8 string = "B8"
	BAMBOO9 string = "B9"

	//VALUE 1-9
	CHARACTER1 string = "C1"
	CHARACTER2 string = "C2"
	CHARACTER3 string = "C3"
	CHARACTER4 string = "C4"
	CHARACTER5 string = "C5"
	CHARACTER6 string = "C6"
	CHARACTER7 string = "C7"
	CHARACTER8 string = "C8"
	CHARACTER9 string = "C9"

	//VALUE 1-9
	DOTS1 string = "D1"
	DOTS2 string = "D2"
	DOTS3 string = "D3"
	DOTS4 string = "D4"
	DOTS5 string = "D5"
	DOTS6 string = "D6"
	DOTS7 string = "D7"
	DOTS8 string = "D8"
	DOTS9 string = "D9"

	//VALUE 1 -7
	EAST  string = "H31"
	SOUTH string = "H42"
	WEST  string = "H53"
	NORTH string = "H64"
	RED   string = "H75"
	GREEN string = "H86"
	WHITE string = "H97"

	//VALUE = 0
	F_PLUMBLOSSOM   string = "F101"
	F_ORCHID        string = "F102"
	F_CHRYSANTHEMUM string = "F103"
	F_BAMBOO        string = "F104"
	F_SPRING        string = "F105"
	F_SUMMER        string = "F106"
	F_AUTUMN        string = "F107"
	F_WINTER        string = "F108"
)

func NewTile(src string) Tile {
	t := Tile{}
	t.From(src)
	return t
}

type Tile struct {
	Suit string
	Rank int8
	Seq  int
}

func (t *Tile) From(src string) {
	t.Suit = src[:1]
	v, _ := strconv.ParseUint(src[1:], 10, 8)
	t.Rank = int8(v)
	t.cn()
}

func (t *Tile) cn() {
	if t.Suit == CHARACTER {
		t.Seq = int(t.Rank)
		return
	}
	if t.Suit == BAMBOO {
		t.Seq = int(10 + t.Rank)
		return
	}
	if t.Suit == DOTS {
		t.Seq = int(20 + t.Rank)
		return
	}
	if t.Suit == HORNOR {
		t.Seq = int(t.Rank)
	}
}

func (t Tile) Name() string {
	return t.Suit + strconv.FormatInt(int64(t.Rank), 10)
}
