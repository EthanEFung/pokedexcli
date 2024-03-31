package main

import (
	"fmt"
	"image"
	_ "image/png"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/qeesung/image2ascii/convert"
)

type detail struct {
	name              string
	no                int
	description       string
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
		d.no = msg.ID
		d.height = decimetresToImperialUnits(msg.Height)
		d.weight = hectogramsToPounds(msg.Weight)
		d.types = []string{}
		for _, t := range msg.Types {
			d.types = append(d.types, t.Type.Name)
		}
	case pokeapi.PokemonSpecies:
		// find the first english entry and use as the description
		for _, entry := range msg.FlavorTextEntries {
			if entry.Language.Name == "en" {
				d.description = strings.Join(strings.Fields(entry.FlavorText), " ")
				break
			}
		}
	case image.Image:
		converter := convert.NewImageConverter()
		options := convert.Options{
			FixedWidth:  64,
			FixedHeight: 32,
			FitScreen:   true,
			Colored:     true,
		}
		d.asciiSprite = converter.Image2ASCIIString(msg, &options)
	}
	return d, nil
}
func (d detail) View() string {
	var s string
	s += fmt.Sprintf("%s #%04d\n", d.name, d.no)
	if d.asciiSprite != "" {
		s += d.asciiSprite + "\n"
	}
	s += d.description + "\n"
	s += "height: " + d.height + "\n"
	s += "weight: " + d.weight + "\n"
	s += "types: "
	for i, t := range d.types {
		s += t
		if i < len(d.types)-1 {
			s += ", "
		}
	}
	s += "\n"

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
