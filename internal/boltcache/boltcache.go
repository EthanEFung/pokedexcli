package boltcache

import "github.com/boltdb/bolt"

type Cache struct {
	db *bolt.DB
}

func NewCache(db *bolt.DB) Cache {
	return Cache{db}
}

func (c Cache) Get(key string) ([]byte, bool) {
	var value []byte
	var found bool
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pokedex"))
		v := b.Get([]byte(key))
		if v != nil {
			value = v
			found = true
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return value, found
}

func (c Cache) Add(key string, value []byte) {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pokedex"))
		err := b.Put([]byte(key), value)
		return err
	})
	if err != nil {
		panic(err)
	}
}
