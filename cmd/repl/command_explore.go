package main

import (
	"errors"
	"fmt"
)

func commandArea(ctg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no location area provided")
	}
	locationAreaName := args[0]
	res, err := ctg.pokeapiClient.GetLocationArea(locationAreaName)
	if err != nil {
		return err
	}
	fmt.Printf("pokemon in %s:\n", locationAreaName)
	for _, encounter := range res.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
