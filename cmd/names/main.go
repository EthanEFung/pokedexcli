package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

func main() {
	client := pokeapi.NewClient(time.Second*10, inmemorycache.NewCache(time.Hour*2))
	f, err := os.Open("./internal/names.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	fmt.Println("Pokemon not found:")
	x := 0
	for scanner.Scan() {
		name := scanner.Text()
		name = strings.ToLower(name)
		name = strings.ReplaceAll(name, "'", "")
		name = strings.ReplaceAll(name, ".", "")
		name = strings.ReplaceAll(name, "♀", "-f")
		name = strings.ReplaceAll(name, "♂", "-m")
		name = strings.ReplaceAll(name, " ", "-")
		name = strings.TrimSpace(name)
		if _, err := client.GetPokemon(name); err != nil {
			fmt.Println(name)
		}
		x++
	}
}
