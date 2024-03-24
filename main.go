package main

import (
	"flag"
	"fmt"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var cacheType string

type config struct {
	pokeapiClient       pokeapi.Client
	caughtPokemon       map[string]pokeapi.Pokemon
	nextLocationURL     *string
	prevLocationURL     *string
	nextLocationAreaURL *string
	prevLocationAreaURL *string
}

func init() {
	flag.StringVar(&cacheType, "cache", "filebased", "the type of cache to use for the application: 'inmemory' or 'filebased'. defaults to 'filebased'")
}

func main() {
	fmt.Println("Hello, World!")
}
