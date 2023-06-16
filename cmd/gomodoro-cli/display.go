package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var numbers = map[int]string{
	0: `
 ██████ 
██    ██
██    ██
██    ██
 ██████ 
`,
	1: `
 ██
███
 ██
 ██
 ██
`,
	2: `
██████ 
     ██
 █████ 
██     
███████
`,
	3: `
██████ 
     ██
 █████ 
     ██
██████
`,
	4: `
██   ██
██   ██
███████
     ██
     ██
`,
	5: `
███████
██     
███████
     ██
███████
`,
	6: `
 ██████ 
██      
███████ 
██    ██
 ██████ 
`,
	7: `
███████
     ██
    ██ 
   ██  
   ██  
`,
	8: `
 █████ 
██   ██
 █████ 
██   ██
 █████ 
`,
	9: `
 █████ 
██   ██
 ██████
     ██
 █████ 
`,
}
var chars = map[rune]string{
	'c': `

██

██

`,
	's': `
  
  
  
  
  
`,
	'f': `
███████  ██████   ██████ ██    ██ ███████
██      ██    ██ ██      ██    ██ ██     
█████   ██    ██ ██      ██    ██ ███████
██      ██    ██ ██      ██    ██      ██
██       ██████   ██████  ██████  ███████
`,
	'p': `
██████   █████  ██    ██ ███████ ███████
██   ██ ██   ██ ██    ██ ██      ██     
██████  ███████ ██    ██ ███████ █████  
██      ██   ██ ██    ██      ██ ██     
██      ██   ██  ██████  ███████ ███████
`,
	'n': `
██ ██████  ██      ███████
██ ██   ██ ██      ██     
██ ██   ██ ██      █████  
██ ██   ██ ██      ██     
██ ██████  ███████ ███████
`,
}

func getNumber(r int) string {
	return numbers[r]
}

func getPhase(p string) string {
	switch p {
	case "Fokusphase", "Focus":
		return focusedStyle.Render(chars['f'])
	case "Verfügbar", "Pause":
		return pausedStyle.Render(chars['p'])
	default:
		return chars['n']
	}
}

func getColonChar() string {
	return chars['c']
}

func getSpaceChar() string {
	return chars['s']
}

func getTimeString(t int64) string {
	m := t / 1000000000 / 60
	s := t / 1000000000 % 60
	m1 := m / 10
	m2 := m % 10
	s1 := s / 10
	s2 := s % 10

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		getNumber(int(m1)),
		getSpaceChar(),
		getNumber(int(m2)),
		getSpaceChar(),
		getColonChar(),
		getSpaceChar(),
		getNumber(int(s1)),
		getSpaceChar(),
		getNumber(int(s2)))
}

func addHelp(t string, help string, height int) string {
	// define a string builder
	var b strings.Builder
	// calculate the number of lines in the help and the timer
	helpLines := strings.Count(help, "\n")
	timerLines := strings.Count(t, "\n")
	// add empty lines between timer and help
	emptyLines := height - helpLines - timerLines - 1
	// if number is negative set it to 0
	if emptyLines < 0 {
		emptyLines = 0
	}
	// add the timer
	b.WriteString(t)
	// add empty lines
	b.WriteString(strings.Repeat("\n", emptyLines))
	// add the help
	b.WriteString(help)

	return b.String()
}

func renderTimer(team Team, remaining int64, name, state string) string {
	var rb strings.Builder

	var lb strings.Builder

	var b strings.Builder

	lb.WriteString("Team Name:\nFocusTimer:\nPauseTimer:\nTimer Status:\n")
	rb.WriteString(
		fmt.Sprintf("\t%s\n\t%d min %d sec\n\t%d min %d sec\n\t%s\n",
			team.Name,
			team.Focus/1000000000/60,
			team.Focus/1000000000%60,
			team.Pause/1000000000/60,
			team.Pause/1000000000%60,
			state),
	)
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lb.String(), rb.String()))
	// Print Time
	b.WriteString(getTimeString(remaining) + "\n")

	// Print Phase
	b.WriteString(getPhase(name) + "\n")

	return b.String()
}
