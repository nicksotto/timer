package main

import "github.com/charmbracelet/lipgloss"

var titleStyle = lipgloss.NewStyle().
	Bold(true).Foreground(lipgloss.Color("#3D5A80")).
	Width(60).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#3D5A80")).
	BorderBottom(true).
	BorderTop(true)

var normalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F0F7EE"))

var hilightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#998650"))

var alertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6666"))

var gradientColors = []string{"#998650", "#5A2A27"}
