package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethanefung/pokedexcli/internal/namefinder"
)

func main() {
	// we want to search the national.json file for all pokemon with forms and print them to the screen
	file, err := os.ReadFile("./internal/namefinder/national.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var generations namefinder.BasicPokemonInfoEntries
	if err := json.Unmarshal(file, &generations); err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}
	for i, generation := range generations {
		fmt.Printf("\n\nGeneration %d\n", i+1)
		for _, pokemon := range generation {
			if pokemon.Form != "" {
				fmt.Printf("ID: %d, Name: %s, %s\n", pokemon.ID, pokemon.Name, pokemon.Form)
			}
		}
	}
}
