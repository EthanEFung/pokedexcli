package filebasedcache

import (
	"bufio"
	"encoding/json"
	"fmt"
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
func (l *Ledger) Restore(c *Cache) error {
	for l.Scan() {
		entry := l.Entry()
		url := entry.Url
		switch entry.Msg {
		case WRITE:
			// just an update?
			if el, ok := c.elements[url]; ok {
				el.Value = entry
				c.list.MoveToFront(el)
				continue
			}

			for c.list.Len() >= c.capacity {
				last := c.list.Back()
				lastEntry := last.Value.(LedgerEntry)
				delete(c.elements, lastEntry.Url)
				c.list.Remove(last)
			}

			c.elements[url] = c.list.PushFront(entry)
		case READ:
			if el, ok := c.elements[url]; ok {
				c.list.MoveToFront(el)
			}
		}
	}
	return nil
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
