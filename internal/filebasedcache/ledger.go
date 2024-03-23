package filebasedcache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
)

type Ledger struct {
	ledgerFile *os.File
	logger     *slog.Logger
	filenames  map[string]struct{}
	scanner    *bufio.Scanner
}

func NewLedger(ledgerFile *os.File) *Ledger {
	scanner := bufio.NewScanner(ledgerFile)
	return &Ledger{
		ledgerFile: ledgerFile,
		logger:     slog.New(slog.NewJSONHandler(ledgerFile, nil)),
		scanner:    scanner,
	}
}

func (l *Ledger) Log(entry LedgerEntry) {
	l.logger.Info(
		string(entry.Msg),
		"filename", entry.Filename,
		"url", entry.Url,
	)
}

func (l *Ledger) Scan() bool {
	return l.scanner.Scan()
}

func (l *Ledger) Text() string {
	return l.scanner.Text()
}

func (l *Ledger) Bytes() []byte {
	return l.scanner.Bytes()
}

// Entry returns the next LedgerEntry in the ledgerFile function panics
// if encountering an error parsing the entry.
func (l *Ledger) Entry() LedgerEntry {
	bytes := l.Bytes()
	entry, err := l.Parse(bytes)
	if err != nil {
		panic(err)
	}
	return entry
}

func (l *Ledger) Read() {
	for l.Scan() {
		fmt.Println(l.Text())
	}
}

// Restore replays the history of the ledger returning the log entries
// the last n writes or reads or that have been persisted (not evicted)

func (l *Ledger) Restore(n int) []LedgerEntry {
	// here what matters is the time that the entry was created
	// and msg, because the msg will determine what entries are returned

	// Idea one
	// the data structures that we would need is
	// a binary tree that represents the history over time
	// and a hashmap that holds reference to the url and the time in which the entry was made

	// writes would traverse the tree and insert the nodes in logN time
	// evicts would look up the time in the hashmap, find the parent, and assign the children
	// of the evicted node
	// reads would first evict the old node and then append the new node.

	// or it could be as simple as having a hashmap of urls to indices
	// and once it's done we loop over all the indices that were created skipping over evicted numbers
	// create another hashmap that uses indices as the key and entries as the

	entryMap := make(map[int]LedgerEntry)
	indices := make(map[string]int)
	index := 0

	for l.Scan() {
		entry := l.Entry()
		switch entry.Msg {
		case WRITE:
			indices[entry.Url] = index
			entryMap[index] = entry
			index++
		case EVICT:
			i, ok := indices[entry.Url]
			if !ok {
				log.Fatalf("attempted to evict a ledger entry not previously created: %+v", entry)
			}
			delete(entryMap, i)
			delete(indices, entry.Url)
		case READ:
			i, ok := indices[entry.Url]
			if !ok {
				log.Fatalf("attempted to read a missing ledger entry: %+v", entry)
			}
			delete(entryMap, i)
			indices[entry.Url] = index // replacing old value
			entryMap[index] = entry
			index++
		}
	}
	entries := []LedgerEntry{}
	for i := 0; i < index; i++ {
		if entry, ok := entryMap[i]; ok {
			entries = append(entries, entry)
		}
	}
	// sliding window algorithm?
	if len(entries) < n {
		return entries
	}
	return entries[len(entries)-n:]
}

func (l *Ledger) Parse(data []byte) (LedgerEntry, error) {
	var entry LedgerEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return LedgerEntry{}, err
	}
	return entry, nil
}

func (l *Ledger) IsUniq(filename string) bool {
	if len(filename) == 0 {
		return false
	}
	return true
}
