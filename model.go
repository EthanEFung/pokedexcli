package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/filebasedcache"
	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

type model struct {
	pokeapiClient  pokeapi.Client
	nextPokemonURL *string
	prevPokemonURL *string
	nextSpeciesURL *string
	prevSpeciesURL *string
	last           string
}

func initialModel(cacheType string) *model {
	fpath, err := filepath.Abs("./internal/filebasedcache/ledger.txt")
	if err != nil {
		fmt.Printf("could not find path to ledger %v", err)
		os.Exit(1)
	}
	dirpath, err := filepath.Abs("./internal/filebasedcache/files")
	ledgerFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("could not open file in write mode %v", err)
		os.Exit(1)
	}
	defer ledgerFile.Close()

	cache := filebasedcache.NewCache(dirpath, ledgerFile)
	if cacheType != "inmemory" && cacheType != "filebased" {
		fmt.Printf("unsupported cache type '%s' specified", cacheType)
		os.Exit(1)
	} else if cacheType == "inmemory" {
		cache = inmemorycache.NewCache(time.Hour * 2)
	}

	return &model{
		pokeapiClient: pokeapi.NewClient(time.Hour, cache),
	}
}

func (c *model) Init() tea.Cmd {
	return nil
}

func (c *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if key.Matches(msg, DefaultKeyMap.Quit) {
			return c, tea.Quit
		}
		c.last = msg.String()

	}

	return c, nil
}

func (c *model) View() string {
	return c.last
}
