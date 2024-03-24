package main

import "fmt"

func commandRelease(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("provide a pokemon name to release")
	}
	pokemonName := args[0]
	if _, ok := cfg.caughtPokemon[pokemonName]; !ok {
		return fmt.Errorf("pokemon %s was not previously caught", pokemonName)
	}
	delete(cfg.caughtPokemon, pokemonName)
	fmt.Printf("released %s\n", pokemonName)
	return nil
}
