package urlshort

import (
	"os"

	"github.com/boltdb/bolt"
)

// BoltDB has an embedded bolt.DB instance
// and is used to implement PairProducer
// having access to the database
// https://stackoverflow.com/questions/28800672/how-to-add-new-methods-to-an-existing-type-in-go
type BoltDB struct {
	*bolt.DB
}

// OpenDB create a BoltDB instance with default options
func OpenBDB(path string, mode os.FileMode) (*BoltDB, error) {
	db, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, err
	}
	return &BoltDB{db}, nil
}

// TODO: Implement logic to get key-vallue pairs
// and convert to array of Pair values
func (bdb *BoltDB) Pair() ([]Pair, error) {
	pairs := []Pair{
		Pair{
			"/fs", "https://www.facebook.com/",
		},
		Pair{
			"/na", "https://na33.salesforce.com/console",
		},
	}
	return pairs, nil
}
