package db

import (
	"github.com/boltdb/bolt"
	"github.com/gisanglee/gicoin/utils"
)

const (
	dbName       string = "blockchain.db"
	dataBucket   string = "data"
	blocksBucket string = "blocks"
	checkpoint   string = "checkpoint"
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

func SaveBlock(hash string, data []byte) {

	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)

		return err
	})

	utils.HandleError(err)

}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})

	utils.HandleError(err)

}

func Checkpoint() []byte {
	var data []byte

	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})

	return data
}
