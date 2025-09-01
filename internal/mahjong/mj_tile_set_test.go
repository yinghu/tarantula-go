package mahjong

import "testing"

func TestFourTileSet(t *testing.T) {
	t4 := NewFourTileSet(1)
	if t4.Eye() {
		t.Errorf("should not be an eye")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	d1 := Tile{Suit: "D", Rank: 1}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	d2 := Tile{Suit: "D", Rank: 2}
	if t4.Append(d1).Allowed(d2) {
		t.Errorf("should not be allowed")
	}
	t4.Append(d1).Append(d1)
	if t4.Full() {
		t.Errorf("should not be full")
	}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	t4.Append(d1)
	if !t4.Full() {
		t.Errorf("should be full")
	}
	ts := t4.Formed().Tiles
	if ts[0] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[1] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[2] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[3] != d1 {
		t.Errorf("Should be %v", d1)
	}
}

func TestThreeTileSet(t *testing.T) {
	t4 := NewThreeTileSet(1)
	if t4.Eye() {
		t.Errorf("should not be an eye")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	d1 := Tile{Suit: "D", Rank: 1}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	d2 := Tile{Suit: "D", Rank: 2}
	if t4.Append(d1).Allowed(d2) {
		t.Errorf("should not be allowed")
	}
	t4.Append(d1)
	if t4.Full() {
		t.Errorf("should not be full")
	}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	t4.Append(d1)
	if !t4.Full() {
		t.Errorf("should be full")
	}
	ts := t4.Formed().Tiles
	if ts[0] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[1] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[2] != d1 {
		t.Errorf("Should be %v", d1)
	}
}

func TestTwoTileSet(t *testing.T) {
	t4 := NewTwoTileSet(1)
	if !t4.Eye() {
		t.Errorf("should be an eye")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	d1 := Tile{Suit: "D", Rank: 1}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	d2 := Tile{Suit: "D", Rank: 2}
	if t4.Append(d1).Allowed(d2) {
		t.Errorf("should not be allowed")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	t4.Append(d1)
	if !t4.Full() {
		t.Errorf("should be full")
	}
	ts := t4.Formed().Tiles
	if ts[0] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[1] != d1 {
		t.Errorf("Should be %v", d1)
	}
}

func TestSequenceTileSet(t *testing.T) {
	t4 := NewSequenceTileSet(1)
	if t4.Eye() {
		t.Errorf("should be an eye")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	d1 := Tile{Suit: "D", Rank: 1}
	if !t4.Allowed(d1) {
		t.Errorf("should be allowed")
	}
	d2 := Tile{Suit: "D", Rank: 2}
	if !t4.Append(d1).Allowed(d2) {
		t.Errorf("should be allowed")
	}
	if t4.Full() {
		t.Errorf("should not be full")
	}
	t4.Append(d2)
	if t4.Allowed(d1) {
		t.Errorf("should not be allowed")
	}
	d3 := Tile{Suit: "D",Rank: 3}
	if !t4.Allowed(d3){
		t.Errorf("should be allowed")
	}
	t4.Append(d3)
	if !t4.Full() {
		t.Errorf("should be full")
	}
	d4 := Tile{Suit: "D",Rank: 3}
	if t4.Allowed(d4){
		t.Errorf("should not be allowed")
	}
	ts := t4.Formed().Tiles
	if ts[0] != d1 {
		t.Errorf("Should be %v", d1)
	}
	if ts[1] != d2 {
		t.Errorf("Should be %v", d2)
	}
	if ts[2] != d3 {
		t.Errorf("Should be %v", d3)
	}
}
