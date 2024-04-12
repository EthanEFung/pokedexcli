package main

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/ethanefung/pokedexcli/internal/soundex"
)

var ErrPokemon = errors.New("error getting pokemon")
var ErrSpecies = errors.New("error getting species")

var encoder = soundex.NewSoundexEncoder()

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

func initializePokelist(client pokeapi.Client) *pokelist {
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Filter = filterFunc
	return &pokelist{
		list:   l,
		items:  items,
		client: client,
	}
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

func filterFunc(term string, items []string) []list.Rank {
	ranks := []list.Rank{}
	code := encoder.Encode(term)

	for i, title := range items {
		icode := encoder.Encode(title)
		if code[0] != icode[0] {
			continue
		}
		complete := true
		for j, c := range code {
			if c == '0' {
				break
			}
			if rune(icode[j]) != c {
				complete = false
				break
			}
		}
		if complete == false {
			continue
		}

		ranks = append(ranks, list.Rank{
			Index:          i,
			MatchedIndexes: []int{},
		})

	}
	return ranks
}

func getPokemonDetails(client pokeapi.Client, name string) tea.Cmd {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "♀", "-f")
	name = strings.ReplaceAll(name, "♂", "-m")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.TrimSpace(name)
	return tea.Batch(
		func() tea.Msg {
			p, err := client.GetPokemon(name)
			if err != nil {
				return ErrPokemon
			}
			return p
		},
		func() tea.Msg {
			ps, err := client.GetPokemonSpecies(name)
			if err != nil {
				return ErrSpecies
			}
			return ps
		})
}
