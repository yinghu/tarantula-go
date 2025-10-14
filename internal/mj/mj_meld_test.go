package mj

import (
	"testing"
)

func TestEyeMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Rank: 1}
	p := []Tile{d1, d1}
	m := Meld{Tiles: p}
	if !m.Eye() {
		t.Errorf("should be an eye")
	}
	//fmt.Printf("Name : %s\n",m.Name())
	m.Tiles = append(m.Tiles, d1)
	if m.Eye() {
		t.Errorf("should not be an eye")
	}
}

func TestChowMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}
	p := []Tile{d1, d2, d3}
	m := Meld{Tiles: p}
	if !m.Chow() {
		t.Errorf("should be a chow")
	}
	
	m.Tiles = append(m.Tiles, d1)
	if m.Chow() {
		t.Errorf("should not be a chow")
	}
	x := []Tile{d1, d1, d2}
	c := Meld{Tiles: x}
	if c.Chow() {
		t.Errorf("should not be a chow")
	}
}

func TestPongMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Rank: 1}
	p := []Tile{d1, d1, d1}
	m := Meld{Tiles: p}
	if !m.Pong() {
		t.Errorf("should be a pong")
	}
	
	m.Tiles = append(m.Tiles, d1)
	if m.Pong() {
		t.Errorf("should not be a pong")
	}

}

func TestKongMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Rank: 1}
	p := []Tile{d1, d1, d1, d1}
	m := Meld{Tiles: p}
	if !m.Kong() {
		t.Errorf("should be a kong")
	}
	m.Tiles = append(m.Tiles, d1)
	if m.Kong() {
		t.Errorf("should not be a kong")
	}

}
