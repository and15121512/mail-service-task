package models

import "time"

type EventType int

const (
	EventCreateType EventType = iota
	EventUpdateType
	EventDeleteType
	EventApproveType
	EventDeclineType
)

type Event struct {
	EventId string
	TaskId  string
	Time    time.Time
	Type    EventType
	Status  int
}
