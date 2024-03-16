package main

import (
	"fmt"
	"sort"
)

func commandHelp(cfg *config, args ...string) error {
	names := make([]string, 0)
	availableCommands := getCommands()
	for name := range availableCommands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		cmd := availableCommands[name]
		fmt.Printf("%s - %s\n", name, cmd.desc)
	}
	fmt.Println("")
	return nil
}
