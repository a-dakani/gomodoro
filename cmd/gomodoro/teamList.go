package main

import "fmt"

type ErrorMsg error

type TimerType string

const (
	Focus TimerType = "focus"
	Pause TimerType = "pause"
)

type Timer struct {
	TimerType TimerType
	Running   bool
	Remaining int64
}

type Team struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Focus int64  `json:"focus"`
	Pause int64  `json:"pause"`
}

func (t Team) FilterValue() string {
	return t.Name
}
func (t Team) Title() string {
	return fmt.Sprintf("Team:%s", t.Slug)
}
func (t Team) Description() string {
	return fmt.Sprintf("F: %d / P: %d", t.Focus/1000000000/60, t.Pause/1000000000/60)
}
