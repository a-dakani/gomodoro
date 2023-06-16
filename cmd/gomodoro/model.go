package main

import (
	"fmt"
	"github.com/a-dakani/gomodoro/pkg/tomodoro"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ErrorMsg error

type sessionState int

const (
	noTeams sessionState = iota
	showList
	showTimer
	showInput
)

type Model struct {
	title          string
	state          sessionState
	sub            chan tomodoro.Message
	ws             *tomodoro.WebSocketClient
	input          textinput.Model
	timerName      string
	timerState     tomodoro.MessageType
	timerRemaining int64
	teamList       list.Model
	help           help.Model
	height         int
	width          int
	err            error
	windowTooSmall bool
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
	tl := list.New([]list.Item{}, delegate, initialListHeight, initialListWidth)
	tl.Title = "Teams"

	return &Model{
		title:          "Gomodoro",
		state:          noTeams,
		sub:            make(chan tomodoro.Message, 100),
		input:          ti,
		timerName:      "Inactive",
		timerState:     tomodoro.TimerStopped,
		timerRemaining: 0,
		teamList:       tl,
		help:           help.New(),
		height:         minimalWindowHeight,
		width:          minimalWindowWidth,
		err:            nil,
		windowTooSmall: false,
	}
}

func (m *Model) Init() tea.Cmd {
	m.loadTeams()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	// if msg type is tea.KeyMsg and is CtrlC, return tea.Quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, Keymap.Quit) {
			return m, tea.Quit
		}
	}
	// if msg type is tea.WindowSizeMsg, set the width and height of the model
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// if the window is too small, set the windowTooSmall flag to true
		if msg.Width < minimalWindowWidth || msg.Height < minimalWindowHeight {
			m.windowTooSmall = true
		} else {
			m.windowTooSmall = false
		}
		m.height = msg.Height - stylesHeight - 1
		m.width = msg.Width - stylesWidth
		m.help.Width = msg.Width - stylesWidth - 2
	}

	// if the window is too small, return the model and don't update
	if m.windowTooSmall {
		return m, nil
	}
	// else check the state of the model
	switch m.state {
	case noTeams:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, Keymap.Add):
				m.state = showInput
			}
		}
	case showInput:
		switch msg := msg.(type) {
		case Team:
			if err := addTeamToFile(msg); err != nil {
				m.err = err
				return m, nil
			}
			// reload the list of teams
			m.loadTeams()
			m.teamList.SetWidth(m.width)
			m.teamList.SetHeight(m.height)
			m.teamList.Help.Width = m.width - stylesWidth
			m.state = showList
			m.teamList, cmd = m.teamList.Update(msg)
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
					m.teamList.SetWidth(m.width)
					m.teamList.SetHeight(m.height)
					m.teamList.Help.Width = m.width - stylesWidth
					m.state = showList
					m.teamList, cmd = m.teamList.Update(msg)
				}
			}
		}
		m.input, cmd = m.input.Update(msg)
	case showList:

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			//case key.Matches(msg, Keymap.Back):
			//	return m, nil
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
				if len(m.teamList.Items()) == 0 {
					m.state = noTeams
				}
			}
		}
		m.teamList.SetWidth(m.width)
		m.teamList.SetHeight(m.height)
		m.teamList.Help.Width = m.width - stylesWidth
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
				m.timerRemaining = 0
				m.timerName = "Inactive"
				m.timerState = tomodoro.TimerStopped
			case tomodoro.Connecting:
				m.timerState = tomodoro.Connecting
			case tomodoro.Connected:
				m.timerState = tomodoro.Connected
			case tomodoro.Terminating:
				m.timerRemaining = 0
				m.timerName = "Inactive"
				m.timerState = tomodoro.Terminating
			case tomodoro.Error:
				m.timerState = tomodoro.Error
				m.err = msg.Error
			}
			return m, m.waitForActivity()
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				m.state = showTimer
			case key.Matches(msg, Keymap.Back):
				m.state = showList
			case key.Matches(msg, Keymap.StartFocus):
				m.err = nil
				return m, m.startFocus()
			case key.Matches(msg, Keymap.StartPause):
				m.err = nil
				return m, m.startPause()
			case key.Matches(msg, Keymap.StopTimer):
				m.err = nil
				return m, m.stopTimer()
			case key.Matches(msg, Keymap.ShowFullHelp):
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}
	return m, cmd
}

func (m *Model) View() string {

	var output string
	output += m.renderTitle()

	if m.windowTooSmall {
		t := fmt.Sprintf("Window too small. Please resize.\n\nMinimum width: %d\nMinimum height: %d", minimalWindowWidth, minimalWindowHeight)
		output += addHelp(t, helpStyle.Render("q/ctrl quit"), m.height)
		return appStyle.Width(m.width).Height(m.height + 1).Render(output)
	}

	if m.err != nil {
		output += m.err.Error() + "\n"
	}
	switch m.state {
	case showList:
		output += m.teamList.View()
	case showTimer:
		t := renderTimer(m.teamList.SelectedItem().(Team), m.timerRemaining, m.timerName, string(m.timerState))
		output += addHelp(t, m.help.View(Keymap), m.height)
	case showInput:
		output += m.input.View()
	case noTeams:
		output += "No teams found. to fetch a team from tomodoro press `+`"
	default:
		output += "Something went wrong. Please try again."
	}

	return appStyle.Width(m.width).Height(m.height).Render(output)
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

func (m *Model) startFocus() tea.Cmd {
	return func() tea.Msg {
		err := startFocus(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return ErrorMsg(err)
		}
		return m.waitForActivity()
	}
}

func (m *Model) startPause() tea.Cmd {
	return func() tea.Msg {
		err := startPause(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return ErrorMsg(err)
		}
		return m.waitForActivity()
	}
}

func (m *Model) stopTimer() tea.Cmd {
	return func() tea.Msg {
		err := stopTimer(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return ErrorMsg(err)
		}
		return m.waitForActivity()
	}
}

func (m *Model) joinTeam() tea.Cmd {
	return func() tea.Msg {
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

func (m *Model) renderTitle() string {
	return titleStyle.Width(m.width-2).Render(fmt.Sprintf("%s", m.title)) + "\n"
}
