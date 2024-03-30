package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

type detail struct {
	name              string
	no                int
	descriptions      []string
	descriptionCursor int
	asciiSprite       string
	height            string
	weight            string
	types             []string
	stats             stats
	evolutions        []string
}

type stats struct {
	hp             string
	attack         string
	defense        string
	specialAttack  string
	specialDefense string
	speed          string
}

func (d detail) Init() tea.Cmd { return nil }
func (d detail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pokeapi.Pokemon:
		d.name = msg.Name
		d.height = decimetresToImperialUnits(msg.Height)
		d.weight = hectogramsToPounds(msg.Weight)
		d.no = msg.ID
	}
	return d, nil
}
func (d detail) View() string {
	var s string
	s += fmt.Sprintf("%s #%04d\n", d.name, d.no)
	s += "height: " + d.height + "\n"
	s += "weight: " + d.weight + "\n"

	return s
}

func decimetresToImperialUnits(decimetres int) string {
	inches := (decimetres * 1000) / 254
	feet := inches / 12
	inches = inches - (feet * 12)
	return fmt.Sprintf("%d'%d\"", feet, inches)
}

func hectogramsToPounds(hectograms int) string {
	var pounds float64 = float64(hectograms) * 0.220462
	return fmt.Sprintf("%.2f lbs", pounds)
}
