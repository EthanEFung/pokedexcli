package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "cursor down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/h", "cursor left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "cursor right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⇥", "next section"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("⇧+⇥", "previous section"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}
