package main

import (
	"github.com/boltdb/bolt"
	"time"
)

func populateDB(pathDB, bucketName string) error {
	db, err := bolt.Open(pathDB, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	pathsToUrls := map[string]string{
		"/bolt-src":       "https://github.com/boltdb/bolt",
		"/bolt-progville": "https://www.progville.com/go/bolt-embedded-db-golang/",
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		for from, to := range pathsToUrls {
			err := bucket.Put([]byte(from), []byte(to))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
