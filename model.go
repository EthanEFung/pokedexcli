package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/pokedexcli/internal/filebasedcache"
	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var docStyle = lipgloss.NewStyle().Margin(1).Border(lipgloss.NormalBorder())

type model struct {
	pokeapiClient  pokeapi.Client
	cursor         int
	nextPokemonURL *string
	prevPokemonURL *string
	nextSpeciesURL *string
	prevSpeciesURL *string
	last           string
	list           list.Model
	detail         tea.Model
	focused        tea.Model
	err            error
}

type listItem struct {
	title string
}

func (li listItem) Title() string       { return li.title }
func (li listItem) Description() string { return "" }
func (li listItem) FilterValue() string { return li.title }

func initialModel(cacheType string) *model {
	fpath, err := filepath.Abs("./internal/filebasedcache/ledger.txt")
	if err != nil {
		fmt.Printf("could not find path to ledger %v", err)
		os.Exit(1)
	}
	dirpath, err := filepath.Abs("./internal/filebasedcache/files")
	ledgerFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o600)

	if err != nil {
		fmt.Printf("could not open file in write mode %v", err)
		os.Exit(1)
	}

	cache := filebasedcache.NewCache(dirpath, ledgerFile)
	if cacheType != "inmemory" && cacheType != "filebased" {
		fmt.Printf("unsupported cache type '%s' specified", cacheType)
		os.Exit(1)
	} else if cacheType == "inmemory" {
		cache = inmemorycache.NewCache(time.Hour * 2)
	}
	client := pokeapi.NewClient(time.Hour, cache)
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	deets := detail{}

	return &model{
		pokeapiClient: client,
		list:          l,
		detail:        deets,
		cursor:        0,
	}
}

func (c *model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			list, err := c.pokeapiClient.GetPokemonList(nil)
			if err != nil {
				return err // feels weird to return different
			}
			return list // types but this is new, so maybe this is just a learning moment
		},
		func() tea.Msg {
			p, err := c.pokeapiClient.GetPokemon("bulbasaur")
			if err != nil {
				return nil
			}
			return p
		},
		func() tea.Msg {
			ps, err := c.pokeapiClient.GetPokemonSpecies("bulbasaur")
			if err != nil {
				return nil
			}
			return ps
		})
}

func (c *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		c.err = msg
	case pokeapi.PokemonList:
		items := []list.Item{}
		for _, p := range msg.Results {
			items = append(items, listItem{p.Name})
		}
		cmd := c.list.SetItems(items)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case pokeapi.Pokemon:
		c.detail, cmd = c.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmd = func() tea.Msg {
			img, err := c.pokeapiClient.GetSprite(msg.Sprites.FrontDefault)
			if err != nil {
				return err
			}
			return img
		}
		cmds = append(cmds, cmd)
	case pokeapi.PokemonSpecies:
		c.detail, cmd = c.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case image.Image:
		// currently the only image that is processed is for the details view
		c.detail, cmd = c.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		c.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return c, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Down):
			c.list.CursorDown()
		case key.Matches(msg, DefaultKeyMap.Up):
			c.list.CursorUp()
		}
	}

	if c.list.Cursor() != c.cursor {
		c.cursor = c.list.Cursor()
		item := c.list.Items()[c.cursor]
		name := item.FilterValue()
		cmd := func() tea.Msg {
			p, err := c.pokeapiClient.GetPokemon(name)
			if err != nil {
				return err
			}
			return p
		}
		cmds = append(cmds, cmd)
		cmd = func() tea.Msg {
			ps, err := c.pokeapiClient.GetPokemonSpecies(name)
			if err != nil {
				return err
			}
			return ps
		}
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c *model) View() string {
	if c.err != nil {
		return fmt.Sprintf("Uh oh, something happened: %+v", c.err)
	}
	listView := docStyle.Render(c.list.View())
	deetsView := docStyle.Render(c.detail.View())
	view := lipgloss.JoinHorizontal(lipgloss.Left, listView, deetsView)

	return view
}
