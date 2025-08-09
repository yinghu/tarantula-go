package main

type Segment struct {
	Id   int64
	Name string
}

type SegementSchedule struct {
	Name     string
	Segments []Segment
}
