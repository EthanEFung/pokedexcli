package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/pokedexcli/internal/namefinder"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var detailsStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
var listStyle = detailsStyle.Copy().Padding(2, 5, 1, 2)

type model struct {
	client pokeapi.Client
	list   tea.Model
	detail tea.Model
	err    error
}

func initialModel(cache pokeapi.Cache) *model {
	client := pokeapi.NewClient(time.Hour, cache)

	pokelist := initializePokelist(client)

	deets := detail{
		client: client,
	}
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
	case namefinder.BasicPokemonInfo:
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
	case tea.WindowSizeMsg:
		m.detail, cmd = m.detail.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.list, cmd = m.list.Update(msg)
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
	listView := listStyle.Render(m.list.View())
	deetsView := detailsStyle.Render(m.detail.View())
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
