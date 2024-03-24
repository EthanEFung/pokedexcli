package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		if len(cleaned) == 0 {
			continue
		}
		commandName := cleaned[0]
		args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}
		availableCommands := getCommands()

		command, ok := availableCommands[commandName]
		if !ok {
			fmt.Println("unknown command")
			continue
		}
		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type cliCommand struct {
	name     string
	desc     string
	callback func(cfg *config, args ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:     "exit",
			desc:     "exits the pokedex",
			callback: commandExit,
		},
		"help": {
			name:     "help",
			desc:     "displays the help menu",
			callback: commandHelp,
		},
		"locations": {
			name:     "locations",
			desc:     "returns a paginated list of locations",
			callback: commandLocations,
		},
		"locationsb": {
			name:     "locationsb",
			desc:     "returns to the previous page of the locations",
			callback: commandLocationsBack,
		},
		"location": {
			name:     "location",
			desc:     "returns the details of a location",
			callback: commandLocation,
		},
		"areas": {
			name:     "areas",
			desc:     "returns a paginated list location areas",
			callback: commandAreas,
		},
		"areasb": {
			name:     "areasb",
			desc:     "returns to the previous page of the location areas",
			callback: commandAreasBack,
		},
		"area": {
			name:     "area",
			desc:     "lists all the pokemon that can be found in a location area",
			callback: commandArea,
		},
		"inspect": {
			name:     "inspect",
			desc:     "displays the details of a pokemon",
			callback: commandInspect,
		},
		"sprite": {
			name:     "sprite",
			desc:     "displays the sprite of a pokemon in ascii art",
			callback: commandSprite,
		},
		"catch": {
			name:     "catch",
			desc:     "adds pokemon to a collection",
			callback: commandCatch,
		},
		"caught": {
			name:     "caught",
			desc:     "lists all caught pokemon",
			callback: commandCaught,
		},
		"release": {
			name:     "release",
			desc:     "removes a pokemon from collection",
			callback: commandRelease,
		},
	}
}

func cleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered)
	return words
}
