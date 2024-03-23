package filebasedcache

import (
	"time"
)

type Msg string

const (
	READ  Msg = "READ"
	WRITE     = "WRITE"
	EVICT     = "EVICT"
)

type LedgerEntry struct {
	Msg      Msg       `json:"msg"`
	Filename string    `json:"filename"`
	Url      string    `json:"url"`
	Time     time.Time `json:"time"`
}
