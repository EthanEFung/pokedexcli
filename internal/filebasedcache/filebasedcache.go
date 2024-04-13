package filebasedcache

import (
	"container/list"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

	"github.com/ethanefung/pokedexcli/internal/pokeapi"
)

var (
	alphanumerics  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ledgerCapacity = 100
)

type Cache struct {
	ledger   *Ledger
	list     *list.List               // doubly linked list
	elements map[string]*list.Element // elements of the list
	mux      *sync.Mutex
	capacity int
	dirpath  string
}

// NewCache returns a filebasedcache that satisfies the interface of the
// pokeapi.Cache. When this is instantiated, under the hood, it will create
// a ledger which holds reference to a file path and a directory path.
func NewCache(dirpath string, ledgerFile *os.File) pokeapi.Cache {
	ledger := NewLedger(ledgerFile)
	c := &Cache{
		dirpath:  dirpath,
		ledger:   ledger,
		capacity: ledgerCapacity,
		list:     list.New(),
		elements: make(map[string]*list.Element),
		mux:      &sync.Mutex{},
	}

	err := ledger.Restore(c)
	if err != nil {
		log.Fatalf("error occured restoring cache from ledger %s", err)
	}
	return c
}

// Add will add the url and data response to the cache. Internally, it interacts with the
// ledger by first checking to see if there is room, and inserting if the cache has
// capacity. Otherwise, it checks the least recently updated (LRU) entry, and removes the lru file.
// At which point, the cache is fine to write to, and the url and data is persisted.
func (c *Cache) Add(url string, data []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	// just an update?
	if el, ok := c.elements[url]; ok {
		prevEntry := el.Value.(LedgerEntry)
		entry := LedgerEntry{
			Msg:      WRITE,
			Url:      prevEntry.Url,
			Filename: prevEntry.Filename,
		}
		c.ledger.Log(entry)
		el.Value = entry

		c.list.MoveToFront(el)
		return
	}

	// insert file
	var filename string
	var isUniqFilename bool
	var attempts int
	// we use the ledgers isUniqFilename because we don't want to add a hash that has been used
	// that already has been written to the ledger
	for isUniqFilename == false && attempts < 20 {
		filename = randSeq(16)
		isUniqFilename = true
		for _, el := range c.elements {
			if filename == el.Value.(LedgerEntry).Filename {
				isUniqFilename = false
				break
			}
		}
		attempts++
	}
	if isUniqFilename == false {
		log.Fatalf("attempted to assign a filename but failed")
	}

	for c.list.Len() >= c.capacity {
		last := c.list.Back()
		lEntry := last.Value.(LedgerEntry)
		lfPath := filepath.Join(c.dirpath, lEntry.Filename)
		if err := os.Remove(lfPath); err != nil {
			log.Fatalf("attempted to remove a file that does not exist %s", lEntry)
		}
		delete(c.elements, last.Value.(LedgerEntry).Url)
		c.list.Remove(last)
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
	c.ledger.Log(writeEntry)

	c.elements[url] = c.list.PushFront(writeEntry)
}

// Get will look for the data and return the data that was stored by the given url.
// This pointer receiver will return boolean indicating if the data was found.
func (c *Cache) Get(url string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if el, ok := c.elements[url]; ok {
		fName := el.Value.(LedgerEntry).Filename
		fPath := filepath.Join(c.dirpath, fName)
		data, err := os.ReadFile(fPath)
		if err != nil {
			log.Fatalf("could not read file %s: %v", fPath, err)
		}
		readEntry := LedgerEntry{
			Msg:      READ,
			Filename: fName,
			Url:      url,
		}
		c.ledger.Log(readEntry)
		c.list.MoveToFront(el)

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
