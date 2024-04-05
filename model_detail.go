package main

import (
	"fmt"
	"image"
	_ "image/png"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/muesli/reflow/wordwrap"
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

type pType string

const (
	normalType   pType = "normal"
	fireType           = "fire"
	waterType          = "water"
	electricType       = "electric"
	grassType          = "grass"
	iceType            = "ice"
	fightingType       = "fighting"
	poisonType         = "poison"
	groundType         = "ground"
	flyingType         = "flying"
	psychicType        = "psychic"
	bugType            = "bug"
	rockType           = "rock"
	ghostType          = "ghost"
	dragonType         = "dragon"
	darkType           = "dark"
	steelType          = "steel"
	fairyType          = "fairy"
)

type pTypeStyle struct {
	pType   pType
	bgColor string
	color   string
}

var pTypeStylesMap = map[string]pTypeStyle{
	"normal":   {normalType, "#A8A77A", "#FFF"},
	"fire":     {fireType, "#EE8130", "#FFF"},
	"water":    {waterType, "#6390F0", "#FFF"},
	"electric": {electricType, "#F7D02C", "#FFF"},
	"grass":    {grassType, "#7AC74C", "#FFF"},
	"ice":      {iceType, "#96D9D6", "#FFF"},
	"fighting": {fightingType, "#C22E28", "#FFF"},
	"poison":   {poisonType, "#A33EA1", "#FFF"},
	"ground":   {groundType, "#E2BF65", "#FFF"},
	"flying":   {flyingType, "#A98FF3", "#FFF"},
	"psychic":  {psychicType, "#F95587", "#FFF"},
	"bug":      {bugType, "#A6B91A", "#FFF"},
	"rock":     {rockType, "#B6A136", "#FFF"},
	"ghost":    {ghostType, "#735797", "#FFF"},
	"dragon":   {dragonType, "#6F35FC", "#FFF"},
	"dark":     {darkType, "#705746", "#FFF"},
	"steel":    {steelType, "#B7B7CE", "#FFF"},
	"fairy":    {fairyType, "#D685AD", "#FFF"},
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
	s += fmt.Sprintf("No #%04d %s\n", d.no, toStartCase(d.name))
	if d.asciiSprite != "" {
		s += d.asciiSprite + "\n"
	}
	s += strings.Repeat("-", 64) + "\n\n"
	s += "types: "
	for _, t := range d.types {
		pts := pTypeStylesMap[t]
		s += pts.String()
	}
	s += "\n"
	s += "height: " + d.height + "\n"
	s += "weight: " + d.weight + "\n\n"
	s += strings.Repeat("-", 64) + "\n\n"
	desc := wordwrap.String(d.description, 62)
	s += desc + "\n"

	return s
}

func (pts pTypeStyle) String() string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pts.bgColor)).
		Foreground(lipgloss.Color(pts.color)).
		Padding(0, 1).
		Render(string(pts.pType))
}

func toStartCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToTitle((s[:1])) + strings.ToLower(s[1:])
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
