package main

import (
	"errors"
	"fmt"
)

func commandLocations(cfg *config, args ...string) error {
	res, err := cfg.pokeapiClient.GetLocations(cfg.nextLocationURL)
	if err != nil {
		return err
	}
	for _, location := range res.Results {
		fmt.Printf(" - %s\n", location.Name)
	}
	cfg.nextLocationURL = res.Next
	cfg.prevLocationURL = res.Previous
	return nil
}

func commandLocationsBack(cfg *config, args ...string) error {
	if cfg.prevLocationURL == nil {
		return errors.New("No previous locations to show")
	}
	res, err := cfg.pokeapiClient.GetLocations(cfg.prevLocationURL)
	if err != nil {
		return err
	}
	for _, location := range res.Results {
		fmt.Printf(" - %s\n", location.Name)
	}
	cfg.nextLocationURL = res.Next
	cfg.prevLocationURL = res.Previous
	return nil
}
