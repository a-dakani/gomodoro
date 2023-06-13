package main

import (
	"fmt"
	"github.com/a-dakani/gomodoro/pkg/tomodoro"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

// TODO handle errors
// TODO make a Timer component
// TODO add logging to file

type sessionState int

type keymap struct {
	Add    key.Binding
	Remove key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Add: key.NewBinding(
		key.WithKeys("+", "a"),
		key.WithHelp("+/a", "add"),
	),
	Remove: key.NewBinding(
		key.WithKeys("-", "d"),
		key.WithHelp("-/d", "remove"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctlr+c"),
		key.WithHelp("q/ctlr+c", "quit"),
	),
}

const (
	noTeams sessionState = iota
	showList
	showTimer
	showInput
)

type Model struct {
	state          sessionState
	sub            chan tomodoro.Message
	ws             *tomodoro.WebSocketClient
	input          textinput.Model
	timerName      string
	timerState     tomodoro.MessageType
	timerRemaining int64
	teamList       list.Model
	err            error
}

func New() *Model {
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Team Slug"
	ti.Focus()

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3)
	delegate.ShortHelpFunc = func() []key.Binding { return []key.Binding{Keymap.Add, Keymap.Remove} }
	delegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{{Keymap.Add, Keymap.Remove}} }
	tl := list.New([]list.Item{}, delegate, defaultListHeight, defaultListWidth)
	tl.Title = "Teams"

	return &Model{
		state:          noTeams,
		sub:            make(chan tomodoro.Message, 100),
		input:          ti,
		timerName:      "None",
		timerState:     tomodoro.TimerStopped,
		timerRemaining: 0,
		teamList:       tl,
		err:            nil,
	}
}

func (m *Model) Init() tea.Cmd {
	m.loadTeams()
	return nil
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
		m.state = showList
	}

	m.teamList.SetItems(items)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// if msg type is tea.KeyMsg and is CtrlC, return tea.Quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// else check the state of the model
	switch m.state {
	case noTeams:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// update the Input width and height
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, Keymap.Add):
				m.state = showInput
			}
		}
	case showInput:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// update the Input width and height
		case Team:
			if err := addTeamToFile(msg); err != nil {
				m.err = err
				return m, nil
			}
			// reload the list of teams
			m.loadTeams()
			m.state = showList
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				// reset the error in case it was set
				m.err = nil
				return m, m.addTeam()
			case key.Matches(msg, Keymap.Back):
				if len(m.teamList.Items()) == 0 {
					m.state = noTeams
				} else {
					m.state = showList
				}
			}
		}
		m.input, cmd = m.input.Update(msg)
	case showList:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// update the list width and height
			m.teamList.SetHeight(msg.Height)
			m.teamList.SetWidth(msg.Width)
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				m.state = showTimer
				return m, tea.Batch(m.joinTeam(), m.waitForActivity())
			case key.Matches(msg, Keymap.Add):
				m.state = showInput
			case key.Matches(msg, Keymap.Remove):

				if err := removeTeamFromFile(m.teamList.Items()[m.teamList.Index()].(Team)); err != nil {
					m.err = err
				}
				m.teamList.RemoveItem(m.teamList.Index())
				//TODO remove team from file
				if len(m.teamList.Items()) == 0 {
					m.state = noTeams
				}
			}

		}
		m.teamList, cmd = m.teamList.Update(msg)
	case showTimer:
		switch msg := msg.(type) {
		case tomodoro.Message:
			switch msg.Type {
			case tomodoro.Tick:
				if m.timerState != tomodoro.TimerStarted {
					m.timerState = tomodoro.TimerStarted
				}
				m.timerName = msg.Payload.Name
				m.timerRemaining = msg.Payload.RemainingTime
			case tomodoro.TimerStarted:
				m.timerState = tomodoro.TimerStarted
				m.timerName = msg.Payload.Name
			case tomodoro.TimerStopped:
				m.timerState = tomodoro.TimerStopped
			case tomodoro.Connecting:
				m.timerState = tomodoro.Connecting
			case tomodoro.Connected:
				m.timerState = tomodoro.Connected
			case tomodoro.Terminating:
				m.timerRemaining = 0
				m.timerName = "None"
				m.timerState = tomodoro.Terminating
			case tomodoro.Error:
				m.timerState = tomodoro.Error
				m.err = msg.Error

			}
			return m, m.waitForActivity()
		case tea.WindowSizeMsg:
			// update the timer width and height
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				m.state = showTimer
			case key.Matches(msg, Keymap.Back):
				m.state = showList
			}
		}

	}

	return m, cmd
}

func (m *Model) View() string {
	var output string
	if m.err != nil {
		output += m.err.Error() + "\n"
	}
	switch m.state {
	case showList:
		output += m.teamList.View()
	case showTimer:
		output += m.timerString()
	case showInput:
		output += m.input.View()
	case noTeams:
		output += "No teams found. to fetch a team from tomodoro press `+`"
	default:
		output += "Something went wrong. Please try again."
	}

	return output
}

func (m *Model) timerString() string {
	// TODO: build a seperate bubble for the timer in pkg
	var b strings.Builder
	if len(m.teamList.Items()) == 0 {
		b.WriteString("No teams found. Please add a team.")
		return b.String()
	}
	team := m.teamList.SelectedItem().(Team)
	b.WriteString(team.Name + "\n")
	tf := fmt.Sprintf("focusTimer: %d minutes and %d seconds", team.Focus/1000000000/60, team.Focus/1000000000%60)
	tp := fmt.Sprintf("pauseTimer: %d minutes and %d seconds", team.Pause/1000000000/60, team.Pause/1000000000%60)

	b.WriteString(tf + "\n")
	b.WriteString(tp + "\n")
	b.WriteString(fmt.Sprintf("Timer Status:\t %s", m.timerState) + "\n")

	//Print Time
	b.WriteString(getTimeString(m.timerRemaining) + "\n")

	b.WriteString(getPhase(m.timerName) + "\n")

	return b.String()
}

func (m *Model) addTeam() tea.Cmd {
	return func() tea.Msg {
		team, err := getTeam(m.input.Value())
		if err != nil {
			m.err = err
			return ErrorMsg(err)
		}
		return team
	}
}

func (m *Model) joinTeam() tea.Cmd {
	return func() tea.Msg {
		//TODO reset timer and status when joining a new team
		//TODO refactor method to be cancelable
		//TODO start method only when team is selected and timer is shown
		//TODO when another timer is selected, cancel the previous one and connect to the new one
		slug := m.teamList.SelectedItem().(Team).Slug
		// if there is already a websocket connection, check if it is the same team
		if m.ws != nil {
			if m.ws.Slug == slug {
				return nil
			} else {
				m.ws.Stop()
				m.ws = tomodoro.NewWebSocketClient(m.teamList.SelectedItem().(Team).Slug)
				m.ws.Start()
				for {
					for elem := range m.ws.OutChan {
						m.sub <- elem
					}
				}
			}
		}
		m.ws = tomodoro.NewWebSocketClient(m.teamList.SelectedItem().(Team).Slug)
		m.ws.Start()
		for {
			for {
				for elem := range m.ws.OutChan {
					m.sub <- elem
				}
			}
		}
	}
}

func (m *Model) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return <-m.sub
	}
}
