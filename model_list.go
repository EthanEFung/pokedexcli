package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/namefinder"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
	"github.com/ethanefung/pokedexcli/internal/soundex"
)

var ErrPokemon = errors.New("error getting pokemon")
var ErrSpecies = errors.New("error getting species")

var encoder = soundex.NewSoundexEncoder()
var pokemonPath = filepath.Join("./", "internal", "namefinder", "pokemon.json")
var nf = namefinder.NewNameFinder(pokemonPath)

type pokelist struct {
	err    error
	list   list.Model
	items  []list.Item
	client pokeapi.Client
}

type item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Form string `json:"form"`
}

func (i item) Title() string {
	return fmt.Sprintf("#%04d %s", i.ID, i.Name)
}
func (i item) Description() string {
	return i.Form
}
func (i item) FilterValue() string { return i.Name }

func (pl pokelist) Init() tea.Cmd {
	return nil
}

func initializePokelist(client pokeapi.Client) *pokelist {
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "National Dex"
	l.Filter = filterFunc
	l.SetShowPagination(false)
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
	case namefinder.BasicPokemonInfoEntries:
		entries := msg
		for _, gen := range entries {
			for _, pkmn := range gen {
				pl.items = append(pl.items, item{Name: pkmn.Name, ID: pkmn.ID, Form: pkmn.Form})
			}
		}
		cmd := pl.list.SetItems(pl.items)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		selected := pl.list.SelectedItem().(item)
		info := namefinder.BasicPokemonInfo{ID: selected.ID, Name: selected.Name, Form: selected.Form}

		cmd = getPokemonBasicInfo(pl.client, info)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		h, v := detailsStyle.GetFrameSize()
		pl.list.SetSize(msg.Width-h-1, msg.Height-v-3)
	default:
		pl.list, cmd = pl.list.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		selected, ok := pl.list.SelectedItem().(item)
		info := namefinder.BasicPokemonInfo{ID: selected.ID, Name: selected.Name, Form: selected.Form}
		if ok {
			cmd := getPokemonBasicInfo(pl.client, info)
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
	if len(code) == 0 {
		return ranks
	}

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

func getPokemonBasicInfo(client pokeapi.Client, info namefinder.BasicPokemonInfo) tea.Cmd {
	name := nf.Find(info)
	return tea.Batch(
		func() tea.Msg {
			return info
		},
		func() tea.Msg {
			p, err := client.GetPokemon(name)
			if err != nil {
				return ErrPokemon
			}
			return p
		},
	)
}
