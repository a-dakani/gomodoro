package main

import "fmt"

// Team is a struct to hold the team data for the team list widget.
type Team struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Focus int64  `json:"focus"`
	Pause int64  `json:"pause"`
}

// Methods are used to implement the interface for the list widget

// FilterValue returns the value to filter the list by
func (t Team) FilterValue() string {
	return t.Name
}

// Title returns the title of the list item
func (t Team) Title() string {
	return t.Slug
}

// Description returns the description of the list item
func (t Team) Description() string {
	return fmt.Sprintf("Focus: %d min\nPause: %d min", t.Focus/1000000000/60, t.Pause/1000000000/60)
}
