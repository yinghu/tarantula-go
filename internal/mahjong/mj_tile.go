package mahjong

const (
	B1 string = "B1"
	B2 string = "B2"
	B3 string = "B3"
	B4 string = "B4"
	B5 string = "B5"
	B6 string = "B6"
	B7 string = "B7"
	B8 string = "B8"
	B9 string = "B9"
	
	C1 string = "C1"
	C2 string = "C2"
	C3 string = "C3"
	C4 string = "C4"
	C5 string = "C5"
	C6 string = "C6"
	C7 string = "C7"
	C8 string = "C8"
	C9 string = "C9"
	
	D1 string = "D1"
	D2 string = "D2"
	D3 string = "D3"
	D4 string = "D4"
	D5 string = "D5"
	D6 string = "D6"
	D7 string = "D7"
	D8 string = "D8"
	D9 string = "D9"

	H1 string = "H1"
	H2 string = "H2"
	H3 string = "H3"
	H4 string = "H4"
	H5 string = "H5"
	H6 string = "H6"
	H7 string = "H7"

	F1 string = "F1"
	F2 string = "F2"
	F3 string = "F3"
	F4 string = "F4"
	F5 string = "F5"
	F6 string = "F6"
	F7 string = "F7"
	F8 string = "F8"
)

type Tile struct {
	Suit  string
	Value uint8
}

type Stack struct {
	Tiles []Tile
}

func (stack *Stack) Shuffle() {
	
}
