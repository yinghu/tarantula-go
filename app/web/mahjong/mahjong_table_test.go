package main

import (
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/mj"
)

func TestTableSetup(t *testing.T) {
	mt := MahjongTable{}
	mt.New()
	st, err := mt.Sit(100)
	if err != nil {
		t.Errorf("shoud not be error %s", err.Error())
	}
	if st != 0 {
		t.Errorf("Seat should be 0 %d\n", st)
	}
	st, err = mt.Sit(200)
	if err != nil {
		t.Errorf("shoud not be error %s", err.Error())
	}
	if st != 1 {
		t.Errorf("Seat should be 1 %d\n", st)
	}
	st, err = mt.Sit(300)
	if err != nil {
		t.Errorf("shoud not be error %s", err.Error())
	}
	if st != 2 {
		t.Errorf("Seat should be 2 %d\n", st)
	}
	st, err = mt.Sit(400)
	if err != nil {
		t.Errorf("shoud not be error %s", err.Error())
	}
	if st != 3 {
		t.Errorf("Seat should be 3 %d\n", st)
	}
	mt.Dice()
	dealer, err := mt.Deal()
	if err != nil {
		t.Errorf("shoud not be error %s", err.Error())
	}
	d := (mt.Pts() - 1) % 4
	if dealer != d {
		t.Errorf("dealer should be %d %d\n", dealer, d)
	}
	mp := mt.Players[d]
	hz := mp.Hand.TileSize()
	if hz != 14 {
		t.Errorf("dealer HAND should be 14 %d\n", hz)
	}
	if !mp.TN {
		t.Errorf("dealer tn should be true %v\n", mp.TN)
	}
	var tc [4]int
	//fz :=0
	for x := range mp.Hand.Tiles {
		tile := mp.Hand.Tiles[x]
		if tile.Suit == mj.CHARACTER {
			tc[TC_C]++
		}
		if tile.Suit == mj.BAMBOO {
			tc[TC_B]++
		}
		if tile.Suit == mj.DOTS {
			tc[TC_D]++
		}
		if tile.Suit == mj.HORNOR {
			tc[TC_H]++
		}
		//if tile.Suit == mj.FLOWER{
		//fz++
		//}
	}

	if tc != mp.TC {
		t.Errorf("should be same numbers %v == %v", tc, mp.TC)
	}
	if len(mp.C) != tc[TC_C] {
		t.Errorf("Character number should not be %d : %d", len(mp.C), tc[TC_C])
	}
	if len(mp.D) != tc[TC_D] {
		t.Errorf("dots number should not be %d : %d", len(mp.D), tc[TC_D])
	}
	if len(mp.B) != tc[TC_B] {
		t.Errorf("bamboo number should not be %d : %d", len(mp.B), tc[TC_B])
	}
	hs := len(mp.HE) + len(mp.HS) + len(mp.HW) + len(mp.HN) + len(mp.R) + len(mp.G) + len(mp.W)
	if hs != tc[TC_H] {
		t.Errorf("hornor number should not be %d : %d", hs, tc[TC_H])

	}
	if mp.LD == 0 {
		t.Errorf("last draw should not be 0 %d", mp.LD)
	}

	for i := range 4 {
		if i != d {
			p := mt.Players[i]
			pz := p.Hand.TileSize()
			if pz != 13 {
				t.Errorf("player HAND should be 13 %d\n", pz)
			}
			if p.TN {
				t.Errorf("player tn should be false %v\n", p.TN)
			}
			var pc [4]int
			for x := range p.Hand.Tiles {
				tile := p.Hand.Tiles[x]
				if tile.Suit == mj.CHARACTER {
					pc[TC_C]++
				}
				if tile.Suit == mj.BAMBOO {
					pc[TC_B]++
				}
				if tile.Suit == mj.DOTS {
					pc[TC_D]++
				}
				if tile.Suit == mj.HORNOR {
					pc[TC_H]++
				}
			}
			if pc != p.TC {
				t.Errorf("should be same numbers %v == %v", pc, p.TC)
			}
			if len(p.C) != pc[TC_C] {
				t.Errorf("Character number should not be %d : %d", len(p.C), pc[TC_C])
			}
			if len(p.D) != pc[TC_D] {
				t.Errorf("dots number should not be %d : %d", len(p.D), pc[TC_D])
			}
			if len(p.B) != pc[TC_B] {
				t.Errorf("bamboo number should not be %d : %d", len(p.B), pc[TC_B])
			}
			hs := len(p.HE) + len(p.HS) + len(p.HW) + len(p.HN) + len(p.R) + len(p.G) + len(p.W)
			if hs != pc[TC_H] {
				t.Errorf("hornor number should not be %d : %d", hs, pc[TC_H])

			}
			if p.LD == 0 {
				t.Errorf("last draw should not be 0 %d", p.LD)
			}
		}
	}

}

func TestTablePlayDrawDiscard(t *testing.T) {
	core.CreateTestLog()
	mt := MahjongTable{}
	mt.New()
	mt.Sit(100)
	mt.Dice()
	mt.Deal()
	dealer := (mt.Pts() - 1) % 4
	if mt.Players[dealer].TileSize() != 14 {
		t.Errorf("dealer hand should be 14 %d", mt.Players[dealer].TileSize())
	}
	if !mt.Players[dealer].TN {
		t.Errorf("dealer tn should be true %v", mt.Players[dealer].TN)
	}
	for i := range 4 {
		if i != dealer {
			del := mt.Players[i].autoPick()
			err := mt.Discard(i, del.Seq)
			if err == nil {
				t.Errorf("should not allowed to discard from tn %v", mt.Players[i].TN)
			}
			if mt.Players[i].TileSize() != 13 {
				t.Errorf("player hand size should be still 13 %d", mt.Players[i].TileSize())
			}
			err = mt.Draw(i)
			if err != nil {
				t.Errorf("player should be able to draw %s", err.Error())
			}
			if mt.Players[i].TileSize() != 14 {
				t.Errorf("player hand size should be still 14 %d", mt.Players[i].TileSize())
			}
			err = mt.Discard(i, del.Seq)
			if err != nil {
				t.Errorf("should be allowed to discard from tn %v", mt.Players[i].TN)
			}
			if mt.Players[i].TileSize() != 13 {
				t.Errorf("player hand size should be still 13 %d", mt.Players[i].TileSize())
			}
		}
	}
	del := mt.Players[dealer].autoPick()
	err := mt.Discard(dealer, del.Seq)
	if err != nil {
		t.Errorf("dealer should be allowed to discard from %v", mt.Players[dealer].TN)
	}
	if mt.Players[dealer].TileSize() != 13 {
		t.Errorf("dealer hand size should be still 13 %d", mt.Players[dealer].TileSize())
	}
	err = mt.Draw(dealer)
	if err != nil {
		t.Errorf("dealer should be allowed to draw from %v", mt.Players[dealer].TN)
	}
	if mt.Players[dealer].TileSize() != 14 {
		t.Errorf("dealer hand size should be still 14 %d", mt.Players[dealer].TileSize())
	}
}

func TestTablePlayChow(t *testing.T) {
	core.CreateTestLog()
	mt := MahjongTable{}
	mt.New()
	mt.Sit(100)
	mt.Dice()
	dealer := (mt.Pts() - 1) % 4
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO2), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO6), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.GREEN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.F_AUTUMN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].TN = true
	err := mt.Draw(dealer)
	if err == nil {
		t.Errorf("draw not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Chow(dealer, mj.FromS(mj.CHARACTER4).Seq, mj.FromS(mj.CHARACTER3).Seq, mj.FromS(mj.CHARACTER5).Seq)
	if err == nil {
		t.Errorf("chow not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.BAMBOO5).Seq)
	if err == nil {
		t.Errorf("pung not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Kong(dealer, mj.FromS(mj.F_AUTUMN).Seq)
	if err != nil {
		t.Errorf("kong allowed from %v", mt.Players[dealer].TN)
	}
	if !mt.Players[dealer].TN {
		t.Errorf("ttn should still true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.GREEN).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Chow(dealer, mj.FromS(mj.CHARACTER4).Seq, mj.FromS(mj.CHARACTER3).Seq, mj.FromS(mj.CHARACTER5).Seq)
	if err != nil {
		t.Errorf("chow allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	if !mt.Players[dealer].TN {
		t.Errorf("tn should be true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.WHITE).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Chow(dealer, mj.FromS(mj.BAMBOO4).Seq, mj.FromS(mj.BAMBOO5).Seq, mj.FromS(mj.BAMBOO6).Seq)
	if err != nil {
		t.Errorf("chow allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	err = mt.Discard(dealer, mj.FromS(mj.WHITE).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}

	for i := range mt.Players[dealer].Formed {
		if mt.Players[dealer].Formed[i].Name() != "C3.C4.C5" && mt.Players[dealer].Formed[i].Name() != "B4.B5.B6" {
			t.Errorf("meld should be C3.C4.C5 or B4.B5.B6 %s", mt.Players[dealer].Formed[i].Name())
		}
	}
	err = mt.Draw(dealer)
	if err != nil {
		t.Errorf("draw allowed from %v", mt.Players[dealer].TN)
	}
}

func TestTablePlayPung(t *testing.T) {
	core.CreateTestLog()
	mt := MahjongTable{}
	mt.New()
	mt.Sit(100)
	mt.Dice()
	dealer := (mt.Pts() - 1) % 4
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO2), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO6), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.GREEN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.F_AUTUMN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].TN = true
	err := mt.Draw(dealer)
	if err == nil {
		t.Errorf("draw not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Chow(dealer, mj.FromS(mj.CHARACTER4).Seq, mj.FromS(mj.CHARACTER3).Seq, mj.FromS(mj.CHARACTER5).Seq)
	if err == nil {
		t.Errorf("chow not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.BAMBOO5).Seq)
	if err == nil {
		t.Errorf("pung not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Kong(dealer, mj.FromS(mj.F_AUTUMN).Seq)
	if err != nil {
		t.Errorf("kong allowed from %v", mt.Players[dealer].TN)
	}
	if !mt.Players[dealer].TN {
		t.Errorf("ttn should still true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.GREEN).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.BAMBOO5).Seq)
	if err != nil {
		t.Errorf("pung allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	if !mt.Players[dealer].TN {
		t.Errorf("tn should be true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.BAMBOO1).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	for i := range mt.Players[dealer].Formed {
		if mt.Players[dealer].Formed[i].Name() != "B5.B5.B5" {
			t.Errorf("meld should be B5.B5.B5 %s", mt.Players[dealer].Formed[i].Name())
		}
	}
	err = mt.Draw(dealer)
	if err != nil {
		t.Errorf("draw allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.BAMBOO2).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.RED).Seq)
	if err != nil {
		t.Errorf("pung allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	for i := range mt.Players[dealer].Formed {
		if mt.Players[dealer].Formed[i].Name() != "B5.B5.B5" && mt.Players[dealer].Formed[i].Name() != "H39.H39.H39" {
			t.Errorf("meld should be B5.B5.B5 OR H39.H39.H39 %s", mt.Players[dealer].Formed[i].Name())
		}
	}
}

func TestTablePlayKong(t *testing.T) {
	core.CreateTestLog()
	mt := MahjongTable{}
	mt.New()
	mt.Sit(100)
	mt.Dice()
	dealer := (mt.Pts() - 1) % 4
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO2), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.BAMBOO6), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER1), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER3), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.CHARACTER5), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.RED), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.GREEN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.F_AUTUMN), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].AppendForTest(mj.FromS(mj.WHITE), false)
	mt.Players[dealer].TN = true
	err := mt.Draw(dealer)
	if err == nil {
		t.Errorf("draw not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Chow(dealer, mj.FromS(mj.CHARACTER4).Seq, mj.FromS(mj.CHARACTER3).Seq, mj.FromS(mj.CHARACTER5).Seq)
	if err == nil {
		t.Errorf("chow not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.BAMBOO5).Seq)
	if err == nil {
		t.Errorf("pung not allow from %v", mt.Players[dealer].TN)
	}
	err = mt.Kong(dealer, mj.FromS(mj.F_AUTUMN).Seq)
	if err != nil {
		t.Errorf("kong allowed from %v", mt.Players[dealer].TN)
	}
	if !mt.Players[dealer].TN {
		t.Errorf("ttn should still true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.GREEN).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.BAMBOO5).Seq)
	if err != nil {
		t.Errorf("pung allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	if !mt.Players[dealer].TN {
		t.Errorf("tn should be true %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.BAMBOO1).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	for i := range mt.Players[dealer].Formed {
		if mt.Players[dealer].Formed[i].Name() != "B5.B5.B5" {
			t.Errorf("meld should be B5.B5.B5 %s", mt.Players[dealer].Formed[i].Name())
		}
	}
	err = mt.Draw(dealer)
	if err != nil {
		t.Errorf("draw allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Discard(dealer, mj.FromS(mj.BAMBOO2).Seq)
	if err != nil {
		t.Errorf("discard allowed from %v", mt.Players[dealer].TN)
	}
	err = mt.Pung(dealer, mj.FromS(mj.RED).Seq)
	if err != nil {
		t.Errorf("pung allow from %v : %s", mt.Players[dealer].TN, err.Error())
	}
	for i := range mt.Players[dealer].Formed {
		if mt.Players[dealer].Formed[i].Name() != "B5.B5.B5" && mt.Players[dealer].Formed[i].Name() != "H39.H39.H39" {
			t.Errorf("meld should be B5.B5.B5 OR H39.H39.H39 %s", mt.Players[dealer].Formed[i].Name())
		}
	}
}
