package namefinder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

func TestNameResolver(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %v", err)
	}
	nationalPath := filepath.Join(wd, "national.json")
	f, err := os.ReadFile(nationalPath)
	if err != nil {
		t.Fatalf("Error reading file %s: %v", "./national.json", err)
	}

	var generations BasicPokemonInfoEntries
	if err := json.Unmarshal(f, &generations); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}
	allPokemon := []BasicPokemonInfo{}
	for _, generation := range generations {
		for _, pokemon := range generation {
			allPokemon = append(allPokemon, pokemon)
		}
	}
	pokemonPath := filepath.Join(wd, "pokemon.json")
	f, err = os.ReadFile(pokemonPath)
	if err != nil {
		t.Fatalf("Error reading file %s: %v", "./pokemon.json", err)
	}
	var response pokeapi.PokemonList
	if err := json.Unmarshal(f, &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}
	keys := make(map[string]struct{})
	for _, pokemon := range response.Results {
		keys[pokemon.Name] = struct{}{}
	}

	nf := NewNameFinder(pokemonPath)
	for _, pokemon := range allPokemon {
		key := nf.Find(pokemon)

		if _, ok := keys[key]; !ok {
			t.Fatalf("Expected key for %s not found: %+v key: %s", pokemon.Name, pokemon, key)
		} else {
			t.Logf("Found key for %s: %s", pokemon.Name, key)
		}
	}
}
