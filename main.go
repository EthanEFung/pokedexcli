package main

import (
	"flag"
	"fmt"
	"log"
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
			log.Fatalf("could not find path to ledger %v", err)
		}
		dirpath, err := filepath.Abs("./internal/filebasedcache/files")
		if err != nil {
			log.Fatalf("could not find path to files %v", err)
		}
		ledgerFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o600)
		if err != nil {
			log.Fatalf("could not open file in write mode %v", err)
		}
		defer ledgerFile.Close()
		cache = filebasedcache.NewCache(dirpath, ledgerFile)
	case "bolt":
		fpath, err := filepath.Abs("./internal/boltcache/pokedex.db")
		if err != nil {
			log.Fatalf("could not find path to bolt db %v", err)
		}
		db, err := bolt.Open(fpath, 0600, nil)
		if err != nil {
			log.Fatalf("could not open bolt db %v", err)
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
			log.Fatalf("could not update bucket %v", err)
		}
		cache = boltcache.NewCache(db)
	default:
		log.Fatalf("unsupported cache type '%s' specified", cacheType)
	}

	p := tea.NewProgram(initialModel(cache), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not run program %v", err)
	}
}
