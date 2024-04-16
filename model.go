package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/pokedexcli/internal/filebasedcache"
	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/namefinder"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var docStyle = lipgloss.NewStyle().Margin(1).Border(lipgloss.NormalBorder())

type model struct {
	client pokeapi.Client
	list   tea.Model
	detail tea.Model
	err    error
}

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

	pokelist := initializePokelist(client)

	deets := detail{}
	return &model{
		client: client,
		list:   pokelist,
		detail: deets,
	}
}

func (m *model) Init() tea.Cmd {
	return readPokemonList()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		m.detail, cmd = m.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case pokeapi.Pokemon:
		m.detail, cmd = m.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmd = getPokemonSprite(m.client, msg.Sprites.FrontDefault)
		cmds = append(cmds, cmd)
	case pokeapi.PokemonSpecies:
		m.detail, cmd = m.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case image.Image:
		// currently the only image that is processed is for the details view
		m.detail, cmd = m.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	default:
		m.list, cmd = m.list.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Uh oh, something happened: %+v", m.err)
	}
	listView := docStyle.Render(m.list.View())
	deetsView := docStyle.Render(m.detail.View())
	view := lipgloss.JoinHorizontal(lipgloss.Left, listView, deetsView)

	return view
}

func readPokemonList() tea.Cmd {
	return func() tea.Msg {
		b, err := os.ReadFile("./internal/namefinder/national.json")
		if err != nil {
			return err
		}
		var entries namefinder.BasicPokemonInfoEntries
		if err := json.Unmarshal(b, &entries); err != nil {
			return err
		}
		return entries
	}
}

func getPokemonSprite(client pokeapi.Client, url string) tea.Cmd {
	return func() tea.Msg {
		img, err := client.GetSprite(url)
		if err != nil {
			return err
		}
		return img
	}
}
