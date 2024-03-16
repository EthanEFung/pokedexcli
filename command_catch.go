package main

import "fmt"

func commandCatch(cfg *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing pokemon name")
	}
	pokemonName := args[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}
	cfg.caughtPokemon[pokemon.Name] = pokemon

	fmt.Printf("caught %s\n", pokemon.Name)

	return nil
}
