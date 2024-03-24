package main

import "fmt"

func commandCaught(cfg *config, args ...string) error {
	caught := cfg.caughtPokemon
	fmt.Printf("caught pokemon: %d\n", len(caught))
	for _, p := range caught {
		fmt.Printf(" - %s\n", p.Name)
	}
	return nil
}
