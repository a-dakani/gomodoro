package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	m := New()
	program := tea.NewProgram(m)

	_, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}

}
