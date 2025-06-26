package database

import (
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"urlshort/internal/handlers"
)

const BucketName = "redirects"

func InitDB(path string) *bbolt.DB {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		return err
	})

	if err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}

	return db
}

func SeedIfEmpty(db *bbolt.DB, jsonPath string) error {
	empty, err := isBucketEmpty(db, BucketName)
	if err != nil {
		return fmt.Errorf("checking bucket: %w", err)
	}
	if !empty {
		log.Println("Skipping seeding, DB entries already initialized...")
		return nil // no need to seed
	}

	log.Println("Seeding DB with redirect entries...")
	var entries []handlers.PathUrl
	err = handlers.UnmarshalFile(jsonPath, &entries, handlers.JsonDecoder)
	if err != nil {
		return fmt.Errorf("loading json: %w", err)
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))

		for _, entry := range entries {
			if err := b.Put([]byte(entry.Path), []byte(entry.Url)); err != nil {
				return fmt.Errorf("put entry for path %q: %w", entry.Path, err)
			}
		}
		return nil
	})
}

func isBucketEmpty(db *bbolt.DB, bucketName string) (bool, error) {
	var isEmpty bool
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		// if bucket is nil database is empty
		if b == nil {
			isEmpty = true
			return nil
		}
		c := b.Cursor()
		k, _ := c.First()
		isEmpty = k == nil //check if there is a key and that it is not nil
		return nil
	})

	return isEmpty, err
}
