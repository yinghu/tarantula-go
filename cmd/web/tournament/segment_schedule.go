package main

type Segment struct {
	InstanceId int64  `json:"-"`
	Name       string `json:"Name"`
}

type SegementSchedule struct {
	TournamentId int64     `json:"-"`
	Name         string    `json:"Name"`
	Segments     []Segment `json:"-"`
}
