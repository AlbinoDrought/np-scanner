package notifications

import (
	bolt "go.etcd.io/bbolt"
)

type Guard interface {
	CheckSent(notifiable Notifiable) (bool, error)
	RecordSent(notifiable Notifiable) error
}

var sentNotificationsBucket = []byte("sent-notifications")

type boltGuard struct {
	db *bolt.DB
}

func (bg *boltGuard) CheckSent(notifiable Notifiable) (bool, error) {
	sent := false
	err := bg.db.View(func(tx *bolt.Tx) error {
		value := tx.Bucket(sentNotificationsBucket).Get([]byte(notifiable.ID()))
		sent = value != nil
		return nil
	})
	return sent, err
}

func (bg *boltGuard) RecordSent(notifiable Notifiable) error {
	return bg.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(sentNotificationsBucket).Put([]byte(notifiable.ID()), []byte{1})
	})
}

func newBoltGuard(db *bolt.DB) (Guard, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(sentNotificationsBucket)
		return err
	})
	return &boltGuard{db}, err
}

func OpenBoltGuard(path string) (Guard, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return newBoltGuard(db)
}
