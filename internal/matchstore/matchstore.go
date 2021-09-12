package matchstore

import (
	"encoding/json"

	"go.albinodrought.com/neptunes-pride/internal/matches"
	bolt "go.etcd.io/bbolt"
)

type MatchStore interface {
	SaveMatch(match *matches.Match) error
	FindOrCreateMatch(gameNumber string) (*matches.Match, error)
}

func Open(path string) (MatchStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := boot(db); err != nil {
		return nil, err
	}

	return &boltMatchStore{db}, nil
}

func boot(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("matches")); err != nil {
			return err
		}
		return nil
	})
}

type boltMatchStore struct {
	db *bolt.DB
}

func (store *boltMatchStore) SaveMatch(match *matches.Match) error {
	serialized, err := json.Marshal(match)
	if err != nil {
		return err
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("matches")).Put([]byte(match.GameNumber), serialized)
	})
}

func (store *boltMatchStore) FindOrCreateMatch(gameNumber string) (*matches.Match, error) {
	var foundMatchSerialized []byte

	err := store.db.View(func(tx *bolt.Tx) error {
		foundMatchSerialized = tx.Bucket([]byte("matches")).Get([]byte(gameNumber))
		return nil
	})

	if err != nil {
		return nil, err
	}

	foundMatch := &matches.Match{}

	if foundMatchSerialized == nil {
		foundMatch = matches.NewMatch(gameNumber)
		err = store.SaveMatch(foundMatch)
	} else {
		err = json.Unmarshal(foundMatchSerialized, foundMatch)
	}

	if err != nil {
		return nil, err
	}

	return foundMatch, nil
}
