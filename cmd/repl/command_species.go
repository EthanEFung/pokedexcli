package main

import (
	"errors"
	"fmt"
)

func commandSpecies(c *config, args ...string) error {
	if len(args) < 1 {
		return errors.New("usage: species <name>")
	}

	name := args[0]
	species, err := c.pokeapiClient.GetPokemonSpecies(name)
	if err != nil {
		return err
	}

	fmt.Println("Name:", species.Name)
	fmt.Println("ID:", species.ID)
	fmt.Println("Order:", species.Order)
	fmt.Println(species.FlavorTextEntries[0].FlavorText)
	return nil
}
