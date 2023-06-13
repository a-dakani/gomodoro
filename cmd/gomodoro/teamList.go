package main

import "fmt"

type Team struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Focus int64  `json:"focus"`
	Pause int64  `json:"pause"`
}

// Methods are used to implement the interface for the list widget

func (t Team) FilterValue() string {
	return t.Name
}
func (t Team) Title() string {
	return fmt.Sprintf("%s", t.Slug)
}
func (t Team) Description() string {
	return fmt.Sprintf("Focus: %d min\nPause: %d min", t.Focus/1000000000/60, t.Pause/1000000000/60)
}
