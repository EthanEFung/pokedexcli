package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethanefung/pokedexcli/internal/filebasedcache"
	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var cacheType string

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationURL     *string
	prevLocationURL     *string
	nextLocationAreaURL *string
	prevLocationAreaURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
}

func init() {
	flag.StringVar(&cacheType, "cache", "filebased", "the type of cache to use for the application: 'inmemory' or 'filebased'. defaults to 'filebased'")
}

func main() {
	flag.Parse()
	// cache := inmemorycache.NewCache(time.Hour)
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
	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Hour, cache),
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}
	fmt.Println("Welcome to the Pokedex REPL")
	fmt.Printf("type one of the following commands to get started:\n\n")
	commandHelp(&cfg)
	startRepl(&cfg)
}
