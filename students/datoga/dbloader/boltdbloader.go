package dbloader

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

const DB_FILE = ".my.db"
const BUCKET_NAME = "urls"

type BoltDBLoader struct {
	db *bolt.DB
}

func NewBoltDBLoader() (*BoltDBLoader, error) {
	db, err := bolt.Open(DB_FILE, 0600, nil)

	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &BoltDBLoader{db: db}, nil
}

func (dbloader BoltDBLoader) Close() {
	dbloader.db.Close()
}

func (dbloader BoltDBLoader) ToURLsMap() (map[string]string, error) {

	urls := map[string]string{}

	err := dbloader.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(BUCKET_NAME))

		if b == nil {
			return errors.New("Bucket does not exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			urls[(string)(k)] = string(v)
		}

		return nil
	})

	return urls, err
}

func (dbloader BoltDBLoader) AddURL(shortURL string, longURL string) error {
	return dbloader.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME))
		err := b.Put([]byte(shortURL), []byte(longURL))
		return err
	})
}
