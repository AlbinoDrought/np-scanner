package matchstore

import (
	"encoding/json"
	"errors"
	"strconv"

	"go.albinodrought.com/neptunes-pride/internal/matches"
	"go.albinodrought.com/neptunes-pride/internal/types"
	bolt "go.etcd.io/bbolt"
)

var ErrMatchNotFound = errors.New("match not found")
var ErrSnapshotNotFound = errors.New("snapshot not found")

type MatchStore interface {
	Matches() ([]string, error)
	SaveMatch(match *matches.Match) error
	FindMatchOrFail(gameNumber string) (*matches.Match, error)
	FindOrCreateMatch(gameNumber string) (*matches.Match, error)

	FindSnapshot(gameNumber string, playerID int, time int64) (*types.APIResponse, error)
	SaveSnapshot(gameNumber string, snapshot *types.APIResponse) error
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
		if _, err := tx.CreateBucketIfNotExists([]byte("snapshots")); err != nil {
			return err
		}
		return nil
	})
}

type boltMatchStore struct {
	db *bolt.DB
}

func (store *boltMatchStore) Matches() ([]string, error) {
	gameNumbers := []string{}

	err := store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("matches")).Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			gameNumbers = append(gameNumbers, string(k))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gameNumbers, nil
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

func (store *boltMatchStore) FindMatchOrFail(gameNumber string) (*matches.Match, error) {
	var foundMatchSerialized []byte

	err := store.db.View(func(tx *bolt.Tx) error {
		foundMatchSerialized = tx.Bucket([]byte("matches")).Get([]byte(gameNumber))
		return nil
	})

	if err != nil {
		return nil, err
	}

	if foundMatchSerialized == nil {
		return nil, ErrMatchNotFound
	}

	foundMatch := &matches.Match{}
	err = json.Unmarshal(foundMatchSerialized, foundMatch)

	if err != nil {
		return nil, err
	}

	return foundMatch, nil
}

func (store *boltMatchStore) FindOrCreateMatch(gameNumber string) (*matches.Match, error) {
	foundMatch, err := store.FindMatchOrFail(gameNumber)

	if err == ErrMatchNotFound {
		foundMatch = matches.NewMatch(gameNumber)
		err = store.SaveMatch(foundMatch)
	}

	if err != nil {
		return nil, err
	}

	return foundMatch, nil
}

func (store *boltMatchStore) FindSnapshot(gameNumber string, playerID int, time int64) (*types.APIResponse, error) {
	var foundSnapshotSerialized []byte

	err := store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("snapshots")).Bucket([]byte(gameNumber))
		if bucket == nil {
			return ErrSnapshotNotFound
		}

		bucket = bucket.Bucket([]byte(strconv.Itoa(playerID)))
		if bucket == nil {
			return ErrSnapshotNotFound
		}

		foundSnapshotSerialized = bucket.Get([]byte(strconv.FormatInt(time, 10)))
		if foundSnapshotSerialized == nil {
			return ErrSnapshotNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	foundSnapshot := &types.APIResponse{}
	err = json.Unmarshal(foundSnapshotSerialized, foundSnapshot)
	if err != nil {
		return nil, err
	}

	return foundSnapshot, nil
}

func (store *boltMatchStore) SaveSnapshot(gameNumber string, snapshot *types.APIResponse) error {
	serialized, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.Bucket([]byte("snapshots")).CreateBucketIfNotExists([]byte(gameNumber))
		if err != nil {
			return err
		}

		bucket, err = bucket.CreateBucketIfNotExists([]byte(strconv.Itoa(snapshot.ScanningData.PlayerUID)))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(strconv.FormatInt(snapshot.ScanningData.Now, 10)), serialized)
	})
}
