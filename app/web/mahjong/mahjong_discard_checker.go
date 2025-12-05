package main

import (
	"slices"

	"gameclustering.com/internal/core"
)

type HandSegmenet struct {
	T int
	C int
}

func hmp(a, b HandSegmenet) int {
	return a.C - b.C
}

func CheckDiscard(p *MahjongPlayer) int {
	tcs := make([]HandSegmenet, 0)
	tcs = append(tcs,HandSegmenet{T:TC_H,C: p.TC[TC_H]})
	tcs = append(tcs,HandSegmenet{T:TC_C,C: p.TC[TC_C]})
	tcs = append(tcs,HandSegmenet{T:TC_B,C: p.TC[TC_B]})
	tcs = append(tcs,HandSegmenet{T:TC_D,C: p.TC[TC_D]})
	slices.SortFunc(tcs,hmp)
	if p.TC[TC_H] > 0 {

	}
	if p.TC[TC_C] > 0 {

	}
	if p.TC[TC_B] > 0 {

	}
	if p.TC[TC_D] > 0 {

	}
	core.AppLog.Printf("Distribution %v\n",tcs)
	return 0
}
