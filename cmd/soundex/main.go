/*
main is a package that generates a soundex trie json file of pokemon names
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/ethanefung/pokedexcli/internal/soundex"
)

type entry struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func main() {
	encoder := soundex.NewSoundexEncoder()
	f, err := os.Open("./internal/names.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := bufio.NewScanner(f)
	entries := make([]entry, 0)
	for s.Scan() {
		text, code := s.Text(), encoder.Encode(s.Text())
		entries = append(entries, entry{Name: text, Code: code})
	}
	if err := s.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Code < entries[j].Code
	})

	f, err = os.Create("./internal/soundex/soundex.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b, err := json.Marshal(entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("fin")
}
