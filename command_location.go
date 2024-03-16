package main

import "fmt"

func commandLocation(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location name")
	}
	locationName := args[0]

	res, err := cfg.pokeapiClient.GetLocation(locationName)
	if err != nil {
		return err
	}
	fmt.Printf("%s:\n", res.Name)
	fmt.Printf("region: %s\n", res.Region.Name)
	fmt.Printf("areas:\n")
	for _, area := range res.Areas {
		fmt.Printf(" - %s\n", area.Name)
	}

	return nil
}
