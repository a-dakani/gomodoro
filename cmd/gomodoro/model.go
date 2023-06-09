package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var timerBlockStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.ThickBorder(), false, false, false, true).
	Render

type keyMap struct {
	joinTeam       key.Binding
	getTeam        key.Binding
	removeTeam     key.Binding
	startFocus     key.Binding
	stopTimer      key.Binding
	startPause     key.Binding
	toggleHelpMenu key.Binding
}

func appKeyMap() *keyMap {
	return &keyMap{
		joinTeam: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "join team"),
		),
		getTeam: key.NewBinding(
			key.WithKeys("a", "+"),
			key.WithHelp("a/+", "add team"),
		),
		removeTeam: key.NewBinding(
			key.WithKeys("r", "-"),
			key.WithHelp("r/-", "remove team"),
		),
		startFocus: key.NewBinding(
			key.WithKeys("ctrl+f"),
			key.WithHelp("ctrl+f", "start focus"),
		),
		startPause: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "start pause"),
		),
		stopTimer: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "stop timer"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
	}
}

type sessionState int

const (
	showTimer sessionState = iota
	showTeams
	noTeams
	showInput
)

type inputFields int

const (
	inputSlugField inputFields = iota
	inputSubmitButton
	inputCancelButton
)

type Model struct {
	state        sessionState
	input        textinput.Model
	timer        timer.Model
	teamList     list.Model
	selectedTeam int
	err          error
}

func New() *Model {
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Team Slug"
	ti.Focus()
	delegate := list.NewDefaultDelegate()
	tl := list.New([]list.Item{}, delegate, defaultListWidth, defaultListWidth)
	tl.Title = "Teams"

	return &Model{
		state:        noTeams,
		input:        ti,
		timer:        timer.Model{},
		teamList:     tl,
		selectedTeam: 0,
		err:          nil,
	}

}

func (m *Model) Init() tea.Cmd {
	m.loadTeams()
	return m.timer.Init()
}

func (m *Model) loadTeams() {
	var teams []Team
	teams, err := readTeamsFile()
	if err != nil {
		teams = []Team{}
		m.err = err
	}
	items := make([]list.Item, len(teams))
	for i, team := range teams {
		items[i] = team
	}
	if len(items) != 0 {
		m.state = showTeams
	}

	m.teamList.SetItems(items)
}

func (m *Model) reScaleList(width, height int) {
	m.teamList.SetWidth(width)
	m.teamList.SetHeight(height)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch ms := msg.(type) {
	case tea.WindowSizeMsg:
		m.reScaleList(ms.Width/3, ms.Height)
	case tea.KeyMsg:
		switch ms.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			m.selectedTeam = m.teamList.Index()
			m.state = showTimer
		case tea.KeyCtrlB:
			m.state = showTeams
		}
	}
	m.teamList, cmd = m.teamList.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	switch m.state {
	case showTeams:
		return m.teamList.View()
	case showTimer:
		return m.timerString()
	case showInput:
		return m.input.View()
	}
	return "No teams found. Please add a team."
}

func (m *Model) timerString() string {
	var b strings.Builder
	if len(m.teamList.Items()) == 0 {
		b.WriteString("No teams found. Please add a team.")
		return timerBlockStyle(b.String())
	}
	team := m.teamList.Items()[m.selectedTeam].(Team)
	b.WriteString(team.Name + "\n")
	tf := fmt.Sprintf("Focus: %d minutes and %d seconds", team.Focus/60/1000000000, team.Focus%60/1000000000)
	tp := fmt.Sprintf("Pause: %d minutes and %d seconds", team.Pause/60/1000000000, team.Pause%60/1000000000)

	b.WriteString(tf + "\n")
	b.WriteString(tp + "\n")

	b.WriteString("\n\n Countdown \n\n")

	b.WriteString(m.timer.View())

	return timerBlockStyle(b.String())
}
