package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	team      Team
	textInput textinput.Model
	err       error
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your team name/slug name here"
	ti.Focus()
	return Model{
		team:      Team{},
		err:       nil,
		textInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	// initialize the model with a request to the server or load teams from Json
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch ms := msg.(type) {
	case tea.KeyMsg:
		switch ms.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, handleTeamInput(m.textInput.Value())
		}
	case Team:
		m.team = ms
	case ErrorMsg:
		m.err = ms
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	s := m.textInput.View() + "\n\n"
	if m.err != nil {
		s += fmt.Sprintf("Error: %s\n", m.err.Error())
	}
	if !m.team.IsEmpty() {
		s += fmt.Sprintf("Team: %s\n", m.team.Name)
		s += fmt.Sprintf("Slug: %s\n", m.team.Slug)
		s += fmt.Sprintf("Focus: %d minutes and %d seconds\n", m.team.Focus/1000000000/60, m.team.Pause/1000000000%60)
		s += fmt.Sprintf("Pause: %d minutes and %d seconds\n", m.team.Pause/1000000000/60, m.team.Pause/1000000000%60)
	}
	return s
}
