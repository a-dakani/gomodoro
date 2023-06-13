package main

import (
	"github.com/charmbracelet/lipgloss"
)

var numbers = map[int]string{
	0: `
 ██████ 
██  ████
██ ██ ██
████  ██
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
███    ██  ██████  ███    ██ ███████ 
████   ██ ██    ██ ████   ██ ██      
██ ██  ██ ██    ██ ██ ██  ██ █████   
██  ██ ██ ██    ██ ██  ██ ██ ██      
██   ████  ██████  ██   ████ ███████
`,
}

func getNumber(r int) string {
	return numbers[r]
}
func getPauseChar() string {
	return chars['p']
}
func getPhase(p string) string {
	switch p {
	case "Fokusphase":
		return chars['f']
	case "Verfügbar":
		return chars['p']
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
