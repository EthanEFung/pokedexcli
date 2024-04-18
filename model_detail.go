package main

import (
	"fmt"
	"image"
	_ "image/png"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/qeesung/image2ascii/convert"
)

type detail struct {
	err          error
	no           int
	name         string
	description  string
	asciiSprite  string
	height       string
	weight       string
	evolutions   []string
	types        []string
	descriptions []string
	stats        stats
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
	d.err = nil
	switch msg := msg.(type) {
	case error:
		d.err = msg
	case pokeapi.Pokemon:
		d.name = msg.Name
		d.no = msg.ID
		d.height = decimetresToImperialUnits(msg.Height)
		d.weight = hectogramsToPounds(msg.Weight)
		d.types = []string{}
		for _, t := range msg.Types {
			d.types = append(d.types, t.Type.Name)
		}
		for _, s := range msg.Stats {
			switch s.Stat.Name {
			case "hp":
				d.stats.hp = fmt.Sprintf("%d", s.BaseStat)
			case "attack":
				d.stats.attack = fmt.Sprintf("%d", s.BaseStat)
			case "defense":
				d.stats.defense = fmt.Sprintf("%d", s.BaseStat)
			case "special-attack":
				d.stats.specialAttack = fmt.Sprintf("%d", s.BaseStat)
			case "special-defense":
				d.stats.specialDefense = fmt.Sprintf("%d", s.BaseStat)
			case "speed":
				d.stats.speed = fmt.Sprintf("%d", s.BaseStat)
			}
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
	if d.err == ErrPokemon {
		return fmt.Sprintf("Error: %v", d.err)
	}
	var s string
	name := strings.ToLower(d.name)
	name = strings.ReplaceAll(name, "-", " ")
	if name != "" {
		name = strings.ToTitle(string(name[0])) + name[1:]
	}
	s += fmt.Sprintf("\n\n    No #%04d %s\n", d.no, name)
	if d.asciiSprite != "" {
		s += d.asciiSprite + "\n"
	}
	s += d.viewNumbers()
	s += d.viewDescription() + "\n"

	return s
}

func (d detail) viewNumbers() string {
	// var s string
	// lipgloss.
	// 	s += strings.Repeat("-", 64) + "\n\n"
	var left string
	left += "    types: "
	for _, t := range d.types {
		pts := pTypeStylesMap[t]
		left += pts.String()
	}
	left += "\n"
	left += "    height: " + d.height + "\n"
	left += "    weight: " + d.weight + "\n\n"
	left = lipgloss.PlaceHorizontal(34, lipgloss.Left, left)

	var right string
	var statsLeft string
	var statsRight string
	right += "  base stats:\n"
	statsLeft += "  hp : " + d.stats.hp + "\n"
	statsLeft += "  atk: " + d.stats.attack + "\n"
	statsLeft += "  def: " + d.stats.defense + "\n"
	statsRight += " speed  : " + d.stats.speed + "\n"
	statsRight += " spl atk: " + d.stats.specialAttack + "\n"
	statsRight += " spl def: " + d.stats.specialDefense + "\n"
	right += lipgloss.JoinHorizontal(lipgloss.Left, statsLeft, statsRight)
	right = lipgloss.PlaceHorizontal(30, lipgloss.Left, right)

	return lipgloss.JoinHorizontal(lipgloss.Left, left, right) + "\n"
}

func (d detail) viewDescription() string {
	var desc string
	if d.err == ErrSpecies {
		desc = "Error: getting species"
	} else {
		desc = d.description
	}

	return lipgloss.NewStyle().Width(64).Padding(0, 4, 2).Render(desc)
}

func (pts pTypeStyle) String() string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pts.bgColor)).
		Foreground(lipgloss.Color(pts.color)).
		Padding(0, 1).
		Render(string(pts.pType))
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
