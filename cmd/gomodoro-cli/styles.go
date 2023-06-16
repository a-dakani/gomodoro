package main

import "github.com/charmbracelet/lipgloss"

// colors
const (
	appForegroundColor     = "#FFFFFF"
	titleBackgroundColor   = "#FFFFFF"
	titleForegroundColor   = "#000000"
	focusedForegroundColor = "#FF0000"
	pausedForegroundColor  = "#006600"
)

// padding, margin and sizes
const (
	minimalWindowHeight  = 25
	minimalWindowWidth   = 60
	initialListHeight    = 10
	initialListWidth     = 10
	appMarginVertical    = 0
	appMarginHorizontal  = 1
	appPaddingVertical   = 0
	appPaddingHorizontal = 1
	// must be changed if title changes
	titleHeight        = 2
	titlePaddingBottom = 1
	stylesHeight       = appMarginVertical*2 + appPaddingVertical*2 + titlePaddingBottom + titleHeight // 2 for the border
	stylesWidth        = appMarginHorizontal*2 + appPaddingHorizontal*2                                // 2 for the border
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(appPaddingVertical, appPaddingHorizontal).
			Margin(appMarginVertical, appMarginHorizontal).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color(appForegroundColor))
	titleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(titleBackgroundColor)).
			Foreground(lipgloss.Color(titleForegroundColor)).
			Margin(0, 0, titlePaddingBottom, 0)
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(focusedForegroundColor))
	pausedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(pausedForegroundColor))
)
