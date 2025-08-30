package mahjong

import (
	"strconv"
	"strings"
)

const (
	//VALUE 1-9
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
	EAST  string = "H1"
	SOUTH string = "H2"
	WEST  string = "H3"
	NORTH string = "H4"
	RED   string = "H5"
	GREEN string = "H6"
	WHITE string = "H7"

	//VALUE = 0
	F_PLUMBLOSSOM   string = "F1"
	F_ORCHID        string = "F2"
	F_CHRYSANTHEMUM string = "F3"
	F_BAMBOO        string = "F4"
	F_SPRING        string = "F5"
	F_SUMMER        string = "F6"
	F_AUTUMN        string = "F7"
	F_WINTER        string = "F8"
)

type Tile struct {
	Suit string
	Num  int8
}

func (t *Tile) From(src string) {
	if strings.HasPrefix(src, "F") {
		t.Suit = src
		t.Num = 0
		return
	}
	t.Suit = src[:1]
	v, _ := strconv.ParseUint(src[1:2], 10, 8)
	t.Num = int8(v)
}
