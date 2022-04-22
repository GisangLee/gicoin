package db

import (
	"github.com/boltdb/bolt"
	"github.com/gisanglee/gicoin/utils"
)

const (
	dbName       string = "blockchain.db"
	dataBucket   string = "data"
	blocksBucket string = "blocks"
)

var db *bolt.DB

func DB() *bolt.DB {

	if db == nil {
		// initialize
		dbPoinrter, err := bolt.Open(dbName, 0600, nil)
		db = dbPoinrter
		utils.HandleError(err)

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)

			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleError(err)

			return err

		})

		utils.HandleError(err)

	}

	return db
}
