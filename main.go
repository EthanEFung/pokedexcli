package main

import (
	"fmt"
	"time"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationURL     *string
	prevLocationURL     *string
	nextLocationAreaURL *string
	prevLocationAreaURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
}

func main() {
	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}
	fmt.Println("Welcome to the Pokedex REPL")
	fmt.Printf("type one of the following commands to get started:\n\n")
	commandHelp(&cfg)
	startRepl(&cfg)
}
