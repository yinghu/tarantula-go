package main

const (
	CMD_DICE int = 1
	CMD_PLAY int = 2
	CMD_END  int = 3
)

type MahjongPlayToken struct {
	Cmd int
}
