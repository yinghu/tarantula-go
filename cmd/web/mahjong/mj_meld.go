package main

//2 to 4 tiles, or 14 if 13 orphans
const (
	NONE uint8 = 0
	EYE  uint8 = 1
	CHOW uint8 = 2
	PONG uint8 = 3
	KNOG uint8 = 4
)

type Meld struct {
	Tiles []Tile
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

func (m *Meld) Pong() bool {
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
	for i,v := range m.Tiles {
		nm += v.Name()
		if i < sz-1 {
			nm = nm + "."
		}
	}
	return nm
}
