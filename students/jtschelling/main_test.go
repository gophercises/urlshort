package main

import (
  "testing"
  "log"

  "github.com/tidwall/buntdb"
)

func TestSeedDB(t *testing.T) {
  // open in diskless mode
  db, err := buntdb.Open(":memory:")
  if err != nil {
    log.Fatal(err)
  }
  
  defer db.Close()

  seedBuntDB(db)

  err = db.View(func(tx *buntdb.Tx) error {
    _, err := tx.Get("/db")
  	if err != nil{
  		t.Errorf("/db key did not exist")
  	}

    _, err = tx.Get("/db-docs")
    if err != nil{
      t.Errorf("/db-docs key did not exist")
    }
  	return nil
  })
}
