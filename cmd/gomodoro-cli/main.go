package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	m := newModel()
	program := tea.NewProgram(m, tea.WithAltScreen())

	_, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
