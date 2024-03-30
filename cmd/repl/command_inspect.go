package main

import "fmt"

func commandInspect(cfg *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing pokemon name")
	}
	pokemonName := args[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", pokemon.Name)
	fmt.Printf("height: %d\n", pokemon.Height)
	fmt.Printf("weight: %d\n", pokemon.Weight)
	fmt.Printf("types: ")
	for i, t := range pokemon.Types {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(t.Type.Name)
	}
	fmt.Print("\nstats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println()

	return nil
}

func commandList(cfg *config, args ...string) error {
	res, err := cfg.pokeapiClient.GetPokemonList(cfg.nextPokemonListURL)
	if err != nil {
		return err
	}
	for _, pokemon := range res.Results {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	cfg.nextPokemonListURL = res.Next
	cfg.prevPokemonListURL = res.Previous
	return nil
}
