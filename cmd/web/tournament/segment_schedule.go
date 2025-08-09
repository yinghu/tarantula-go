package main

type Segment struct {
	InstanceId   int64
	Name         string
	SegementList []Segment
}

type SegementSchedule struct {
	TournamentId int64 `json:"-"`
	Name         string 
	Segments     []Segment
}
