/*
namefinder is a package that was born out of aggregating data from two sources: bulbagarden and pokeapi.co
Names of all existing pokemon were scraped from bulbagarden, but the api used to get pokemon details come
from pokeapi.co. This space was created to transform names and forms received from bulbagarden into a
uri recognizable by the api.
*/
package namefinder

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"

	"github.com/sahilm/fuzzy"
)

type list []struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (l list) String(i int) string {
	return l[i].Name
}

func (l list) Len() int {
	return len(l)
}

type NameFinder struct {
	list
}

func NewNameFinder(pokemonPath string) *NameFinder {
	b, err := os.ReadFile(pokemonPath)
	if err != nil {
		panic(err)
	}
	var response pokeapi.PokemonList
	if err := json.Unmarshal(b, &response); err != nil {
		panic(err)
	}
	return &NameFinder{
		list: response.Results,
	}
}

// Find takes a basic info and return a dash seperated name recognized as a uri by the pokeapi
func (nf NameFinder) Find(info BasicPokemonInfo) string {
	name := info.Name
	b := strings.Builder{}
	for _, r := range name {
		if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	form := info.Form
	for _, r := range form {
		if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	input := b.String()
	results := []fuzzy.Match{}
	for len(results) == 0 && len(input) > 0 {
		results = fuzzy.FindFrom(input, nf.list)
		input = input[:len(input)-1]
	}

	var best fuzzy.Match
	for _, r := range results {
		if r.Score > best.Score {
			best = r
		}
	}
	return best.Str
}
