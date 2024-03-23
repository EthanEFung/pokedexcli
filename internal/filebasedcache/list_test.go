package filebasedcache

import (
	"testing"
	"time"
)

func TestNewList(t *testing.T) {
	list := NewList(5)
	if list.capacity != 5 {
		t.Errorf("expected list to have a capacity of 5 but got %d", list.capacity)
	}
}

func TestList_Add(t *testing.T) {
	list := NewList(5)
	list.Push(LedgerEntry{})
	if list.capacity != 5 {
		t.Errorf("expected list to have a capacity of 5 but got %d", list.capacity)
	}
}

func TestList_Setup(t *testing.T) {
	entries := []LedgerEntry{{}, {}, {}, {}, {}}
	list := NewList(len(entries))
	list.Setup(entries)
	if list.capacity != 5 {
		t.Errorf("expected list to have a capacity of 5 but got %d", list.capacity)
	}
}

func TestList_Remove(t *testing.T) {

	now := time.Now()
	entries := []LedgerEntry{
		{Filename: "foo", Time: now},
		{Filename: "baz", Time: now},
		{Filename: "foo", Time: now},
		{Filename: "baz", Time: now},
		{}}
	list := NewList(len(entries))
	list.Setup(entries)
	if list.size != 5 {
		t.Errorf("expected list to have a size of 5 but got %d", list.capacity)
	}
	if list.capacity != 5 {
		t.Errorf("expected list to have a capacity of 5 but got %d", list.capacity)
	}

	list.Remove(LedgerEntry{Filename: "foo"})
	list.Reset()
	for list.Scan() {
		t.Logf("%+v", list.Entry())
	}
	if list.size != 3 {
		t.Errorf("expected list to have a size of 3 but got %d", list.capacity)
	}

}
