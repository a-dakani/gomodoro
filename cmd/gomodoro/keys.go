package main

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Add           key.Binding
	Remove        key.Binding
	StartFocus    key.Binding
	StartPause    key.Binding
	StopTimer     key.Binding
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
	Back          key.Binding
	Quit          key.Binding
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
	StartFocus: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "start focus"),
	),
	StartPause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "start pause"),
	),
	StopTimer: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop timer"),
	),
	ShowFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.StartFocus, k.StartPause, k.StopTimer, k.ShowFullHelp}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.StartFocus, k.StartPause, k.StopTimer},
		{k.Back, k.Quit, k.CloseFullHelp},
	}
}
