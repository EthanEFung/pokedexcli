package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanefung/pokedexcli/internal/boltcache"
	"github.com/ethanefung/pokedexcli/internal/filebasedcache"
	"github.com/ethanefung/pokedexcli/internal/inmemorycache"
	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var cacheType string

func init() {
	flag.StringVar(&cacheType, "cache", "bolt", "the type of cache to use for the application: 'inmemory', 'filebased' or 'bolt'. defaults to 'bolt'")
}

func main() {
	flag.Parse()

	var cache pokeapi.Cache
	switch cacheType {
	case "inmemory":
		cache = inmemorycache.NewCache(time.Hour * 2)
	case "filebased":
		fpath, err := filepath.Abs("./internal/filebasedcache/ledger.txt")
		if err != nil {
			fmt.Printf("could not find path to ledger %v", err)
			os.Exit(1)
		}
		dirpath, err := filepath.Abs("./internal/filebasedcache/files")
		ledgerFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o600)
		defer ledgerFile.Close()

		if err != nil {
			fmt.Printf("could not open file in write mode %v", err)
			os.Exit(1)
		}
		cache = filebasedcache.NewCache(dirpath, ledgerFile)
	case "bolt":
		db, err := bolt.Open("./internal/boltcache/pokedex.db", 0600, nil)
		if err != nil {
			fmt.Printf("could not open bolt db %v", err)
			os.Exit(1)
		}
		defer db.Close()
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("pokedex"))
			if err != nil {
				return fmt.Errorf("could not create bucket %v", err)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("could not update bolt db %v", err)
			os.Exit(1)
		}
		cache = boltcache.NewCache(db)
	default:
		fmt.Printf("unsupported cache type '%s' specified", cacheType)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(cache), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
