package urlshort

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// BoltDB has an embedded bolt.DB instance
// and is used to implement PairProducer
// having access to the database
// https://stackoverflow.com/questions/28800672/how-to-add-new-methods-to-an-existing-type-in-go
type BDB struct {
	*bolt.DB
}

// OpenDB create a BoltDB instance with default options
func OpenBDB(path string, mode os.FileMode) (*BDB, error) {
	db, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, err
	}
	return &BDB{db}, nil
}

// LoadInitData is used just to create the pairs bucket and to insert one record.
func (bdb *BDB) LoadInitData() error {
	if err := bdb.Update(func(tx *bolt.Tx) error {
		// create bucket if it doesn't exist
		bk, err := tx.CreateBucketIfNotExists([]byte("pairs"))
		if err != nil {
			return fmt.Errorf("could not create pairs bucket: %v", err)
		}
		// insert one key-value pair
		if err := bk.Put([]byte("/wi"), []byte("https://ru.wikipedia.org")); err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Pair will look for key-value pairs in the "pairs" Bucket
// and will return an array of the Pair structs
func (bdb *BDB) Pair() ([]Pair, error) {
	var pairs []Pair

	if err := bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pairs"))
		b.ForEach(func(k, v []byte) error {
			pairs = append(pairs, Pair{string(k), string(v)})
			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return pairs, nil
}
