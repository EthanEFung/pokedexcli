package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var cacheType string

func init() {
	flag.StringVar(&cacheType, "cache", "filebased", "the type of cache to use for the application: 'inmemory' or 'filebased'. defaults to 'filebased'")
}

func main() {
	flag.Parse()
	p := tea.NewProgram(initialModel(cacheType))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
