package main

import (
  "testing"

  "github.com/tidwall/buntdb"
)

func TestSeedDB(t *testing.T) {
  // open in diskless mode
  db, err := buntdb.Open(":memory:")
  if err != nil {
    log.Fatal(err)
  }

  seedBuntDB(db)

  err := db.View(func(tx *buntdb.Tx) error {
    val, err := tx.Get("/db")
  	if err != nil{
  		t.Errorf("/db key did not exist")
  	}

    val, err := tx.Get("/db-docs")
    if err != nil{
      t.Errorf("/db-docs key did not exist")
    }
  	return nil
  })
}
