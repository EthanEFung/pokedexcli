package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethanefung/pokedexcli/internal/namefinder"
	"github.com/gocolly/colly/v2"
)

func main() {
	list := namefinder.BasicPokemonInfoEntries{}
	c := colly.NewCollector()

	c.OnHTML("table.roundy", func(table *colly.HTMLElement) {
		var id int
		generation := []namefinder.BasicPokemonInfo{}
		table.ForEach("tr[style=\"background:#FFF\"]", func(_ int, tr *colly.HTMLElement) {
			td := tr.DOM.Find("td:has(small)")
			aSel := td.Find("a")
			smallSel := td.Find("small")
			if aSel == nil {
				return
			}
			if smallSel == nil {
				return
			}

			tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
				if strings.HasPrefix(td.Text, "#") {
					no := td.Text[1:]
					newID, err := strconv.Atoi(no)
					if err != nil {
						log.Fatalf("Error converting string to int: %v", err)
						return
					}
					id = newID
				}
			})

			pokemon := namefinder.BasicPokemonInfo{
				ID:   id,
				Name: aSel.Text(),
				Form: smallSel.Text(),
			}
			generation = append(generation, pokemon)
		})
		list = append(list, generation)
	})

	err := c.Visit("https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_National_Pok%C3%A9dex_number")
	if err != nil {
		log.Fatalf("Error visiting page: %v", err)
	}
	b, err := json.Marshal(list)
	if err != nil {
		log.Fatalf("Error marshalling json: %v", err)
	}
	err = os.WriteFile("./internal/namefinder/national.json", b, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("fin")
}
