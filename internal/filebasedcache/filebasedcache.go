package filebasedcache

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var (
	alphanumerics  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ledgerCapacity = 10
)

type Cache struct {
	filenames map[string]string
	ledger    *Ledger
	list      *List
	dirpath   string
	// TODO: should we add a mux here?
}

// NewCache returns a filebasedcache that satisfies the interface of the
// pokeapi.Cache. When this is instantiated, under the hood, it will create
// a ledger which holds reference to a file path and a directory path.
func NewCache(dirpath string, ledgerFile *os.File) pokeapi.Cache {
	ledger := NewLedger(ledgerFile)
	entries := ledger.Restore(ledgerCapacity)
	list := NewList(ledgerCapacity)
	list.Setup(entries)

	return &Cache{
		dirpath:   dirpath,
		ledger:    ledger,
		list:      list,
		filenames: make(map[string]string),
	}
}

// Add will add the url and data response to the cache. Internally, it interacts with the
// ledger by first checking to see if there is room, and inserting if the cache has
// capacity. Otherwise, it checks the least recently updated (LRU) entry, and removes the lru file.
// At which point, the cache is fine to write to, and the url and data is persisted.
func (c *Cache) Add(url string, data []byte) {

	// insert file
	var filename string
	var isUniqFilename bool
	var attempts int
	// we use the ledgers IsUniq method because we don't want to add a hash that has been used
	// that already has been written to the ledger
	for c.ledger.IsUniq(filename) == false && attempts < 20 {
		filename = randSeq(16)
		isUniqFilename = true
		for _, name := range c.filenames {
			if name == filename {
				isUniqFilename = false
				break
			}
		}
		attempts++
	}
	if isUniqFilename == false {
		log.Fatalf("attempted to assign a filename but failed")
	}

	// it's time to write the file
	fPath := filepath.Join(c.dirpath, filename)
	f, err := os.Create(fPath)
	if err != nil {
		log.Fatalf("could not create file with filename %s: %v", filename, err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Fatalf("could not write to file %s: %v", filename, err)
	}

	writeEntry := LedgerEntry{
		Msg:      WRITE,
		Filename: filename,
		Url:      url,
	}
	c.list.Remove(writeEntry)
	c.list.Push(writeEntry)
	c.ledger.Log(writeEntry)

	if c.list.Full() {
		if prev, exists := c.list.Pop(); exists {
			fPath := filepath.Join(c.dirpath, prev.Filename)
			if err := os.Remove(fPath); err != nil {
				log.Fatalf("attempted to remove a file that does not exist")
			}
			entry := LedgerEntry{
				Msg:      EVICT,
				Url:      prev.Url,
				Filename: prev.Filename,
			}
			c.ledger.Log(entry)
		}
	}
}

// Get will look for the data and return the data that was stored by the given url.
// This pointer receiver will return boolean indicating if the data was found.
func (c *Cache) Get(url string) ([]byte, bool) {
	var found bool
	var entry LedgerEntry

	c.list.Reset()
	for c.list.Scan() {
		entry = c.list.Entry()
		if entry.Url == url {
			found = true
			break
		}
	}

	if found {
		fPath := filepath.Join(c.dirpath, entry.Filename)
		data, err := os.ReadFile(fPath)
		if err != nil {
			log.Fatalf("could not read file %s: %v", fPath, err)
		}

		readEntry := LedgerEntry{
			Msg:      READ,
			Filename: entry.Filename,
			Url:      url,
		}
		c.list.Remove(readEntry)
		c.list.Push(readEntry)
		c.ledger.Log(readEntry)

		if c.list.Full() {
			if prev, exists := c.list.Pop(); exists {
				fPath := filepath.Join(c.dirpath, prev.Filename)
				if err := os.Remove(fPath); err != nil {
					log.Fatalf("attempted to remove a file that does not exist")
				}
				evictEntry := LedgerEntry{
					Msg:      EVICT,
					Url:      prev.Url,
					Filename: prev.Filename,
				}
				c.ledger.Log(evictEntry)
			}
		}
		return data, true
	}

	return []byte{}, false
}

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		j := rand.Intn(len(alphanumerics))
		b[i] = alphanumerics[j]
	}
	return string(b)
}
