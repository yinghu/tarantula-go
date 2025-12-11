package mj

import "gameclustering.com/internal/core"

//2 to 4 tiles, or 14 if 13 orphans
const (
	NONE int = 0
	EYE  int = 1
	CHOW int = 2
	PUNG int = 3
	KNOG int = 4
)

type Meld struct {
	Tiles     []Tile `json:"Tiles"`
	Concealed bool
}

func (m *Meld) Type() int {
	if m.Eye() {
		return EYE
	}
	if m.Chow() {
		return CHOW
	}
	if m.Pung() {
		return PUNG
	}
	if m.Kong() {
		return KNOG
	}
	return NONE
}

func (m *Meld) Suit() string {
	return m.Tiles[0].Suit
}

func (m *Meld) Eye() bool {
	if len(m.Tiles) != 2 {
		return false
	}
	return m.Tiles[0] == m.Tiles[1]
}

func (m *Meld) Chow() bool {
	if len(m.Tiles) != 3 {
		return false
	}
	return m.Tiles[0].Suit == m.Tiles[1].Suit && m.Tiles[1].Suit == m.Tiles[2].Suit && m.Tiles[0].Rank+1 == m.Tiles[1].Rank && m.Tiles[1].Rank+1 == m.Tiles[2].Rank
}

func (m *Meld) Pung() bool {
	if len(m.Tiles) != 3 {
		return false
	}
	return m.Tiles[0] == m.Tiles[1] && m.Tiles[1] == m.Tiles[2]

}

func (m *Meld) Kong() bool {
	if len(m.Tiles) != 4 {
		return false
	}
	return m.Tiles[0] == m.Tiles[1] && m.Tiles[1] == m.Tiles[2] && m.Tiles[2] == m.Tiles[3]
}

func (m *Meld) Name() string {
	var nm string
	sz := len(m.Tiles)
	for i, v := range m.Tiles {
		nm += v.Name()
		if i < sz-1 {
			nm = nm + "."
		}
	}
	return nm
}

func (m *Meld) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt32(int32(m.Type())); err != nil {
		return err
	}
	sz := len(m.Tiles)
	if err := buff.WriteInt32(int32(sz)); err != nil {
		return err
	}
	for _, v := range m.Tiles {
		if err := v.Write(buff); err != nil {
			return err
		}
	}
	return nil
}
