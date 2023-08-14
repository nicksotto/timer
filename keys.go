package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	StartWork  key.Binding
	StartBreak key.Binding
	Quit       key.Binding
}

var keys = keyMap{
	StartWork:  key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "start working")),
	StartBreak: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "start a break"), key.WithDisabled()),
	Quit:       key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.StartWork, k.StartBreak, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.StartWork, k.StartBreak, k.Quit},
	}
}
