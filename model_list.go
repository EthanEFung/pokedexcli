package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/ethanefung/pokedexcli/internal/soundex"
)

type pokelist struct {
	err    error
	list   list.Model
	items  []list.Item
	client pokeapi.Client
}

type item struct {
	title string
	code  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

func (pl pokelist) Init() tea.Cmd {
	return nil
}

func (pl pokelist) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		pl.err = msg
	case soundex.Entries:
		entries := msg
		for _, e := range entries {
			pl.items = append(pl.items, item{title: e.Name, code: e.Code})
		}

		cmd := pl.list.SetItems(pl.items)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		selected := pl.list.SelectedItem().(item)
		cmd = getPokemonDetails(pl.client, selected.Title())
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		pl.list.SetSize(msg.Width-h, msg.Height-v)
	default:
		pl.list, cmd = pl.list.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		selected, ok := pl.list.SelectedItem().(item)
		if ok {
			cmd := getPokemonDetails(pl.client, selected.Title())
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}
	return pl, tea.Batch(cmds...)
}

func (fl pokelist) View() string {
	return fl.list.View()
}

func getPokemonDetails(client pokeapi.Client, name string) tea.Cmd {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, "♀", "-f")
	name = strings.ReplaceAll(name, "♂", "-m")
	return tea.Batch(
		func() tea.Msg {
			p, err := client.GetPokemon(name)
			if err != nil {
				return err
			}
			return p
		},
		func() tea.Msg {
			ps, err := client.GetPokemonSpecies(name)
			if err != nil {
				return err
			}
			return ps
		})
}
