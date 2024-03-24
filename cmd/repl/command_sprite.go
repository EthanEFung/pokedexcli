package main

import (
	"errors"
	"fmt"

	"github.com/qeesung/image2ascii/convert"
)

func commandSprite(cfg *config, args ...string) error {
	if len(args) < 1 {
		return errors.New("sprite command requires the name of one pokemon")
	}
	pokemonName := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	spriteImg, err := cfg.pokeapiClient.GetSprite(pokemon.Sprites.FrontDefault)
	if err != nil {
		return err
	}
	converter := convert.NewImageConverter()
	options := convert.Options{
		FixedWidth:  64,
		FixedHeight: 32,
		FitScreen:   true,
		Colored:     true,
	}

	asciiArt := converter.Image2ASCIIString(spriteImg, &options)
	fmt.Println(asciiArt)

	return nil
}
