package main

import (
	"encoding/xml"
	"log"
	"path/filepath"

	"github.com/satori/go.uuid"
	"go.etcd.io/bbolt"
)

// BucketType describes a defined name of a bucket in the database.
type BucketType string

// Enumeration of the database buckets
const (
	DataSetBucket  BucketType = "datasets"
	DataInfoBucket BucketType = "datainfo"
	StockBucket    BucketType = "stocks"
	AccountBucket  BucketType = "accounts"
	CommonBucket   BucketType = "commons"
)

// DB provides the database methods of the application.
type DB struct {
	db *bbolt.DB
}

// InitDB initializes the database in the given folder.
func InitDB(dir string) (*DB, error) {
	path := filepath.Join(dir, "db.data")
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(initDataBuckets)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func initDataBuckets(tx *bbolt.Tx) error {
	buckets := []BucketType{
		DataSetBucket,
		DataInfoBucket,
		StockBucket,
		AccountBucket,
		CommonBucket}
	for _, bucket := range buckets {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
	}
	return nil
}

// Close the database.
func (db *DB) Close() error {
	return db.db.Close()
}

// Get returns the data of the entity with the given ID from the given bucket.
func (db *DB) Get(bucket BucketType, id string) []byte {
	var data []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		log.Println("Error in db.Get: bucket=", bucket, "id=", id)
		return nil
	}
	return data
}

// Put stores the given value under the given id in the bucket.
func (db *DB) Put(bucket BucketType, id string, data []byte) {
	if data == nil {
		return
	}
	err := db.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(id), data)
	})
	if err != nil {
		log.Println("Error while saving data into bucket=",
			bucket, "id=", id, ":>", err)
	}
}

// Delete deletes the entity with the given ID from the bucket.
func (db *DB) Delete(bucket BucketType, id string) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		return bucket.Delete([]byte(id))
	})
}

// Iter iterates over the entries in the given bucket until the function
// parameter returns false or no entries are contained in the database anymore.
func (db *DB) Iter(bucket BucketType, fn func(key, data []byte) bool) {
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if !fn(k, v) {
				break
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error while iterating over", bucket, ":>", err)
	}
}

// RootDataStock returns the root data stock.
func (db *DB) RootDataStock() *DataStock {
	var root *DataStock
	db.Iter(StockBucket, func(key, val []byte) bool {
		ds := &DataStock{}
		xml.Unmarshal(val, ds)
		if ds.IsRoot {
			root = ds
			return false
		}
		return true
	})
	if root != nil {
		return root
	}
	root = &DataStock{
		IsRoot:      true,
		ID:          uuid.NewV4().String(),
		ShortName:   "root",
		Name:        "root",
		Description: "The root data stock"}
	bytes, err := xml.Marshal(root)
	if err != nil {
		log.Println("Error: failed to marshal data stock :>", err)
		return root
	}
	db.Put(StockBucket, root.ID, bytes)
	return root
}

// DataStock returns the data stock with the given ID.
func (db *DB) DataStock(id string) *DataStock {
	var stock *DataStock
	db.Iter(StockBucket, func(key, val []byte) bool {
		if string(key) == id {
			stock = &DataStock{}
			err := xml.Unmarshal(val, stock)
			if err != nil {
				log.Println("Error: loading data stock", id, "failed")
				stock = nil
			}
			return false
		}
		return true
	})
	return stock
}

// Content returns the content of the given data stock.
func (db *DB) Content(stock *DataStock) *InfoList {
	data := db.Get(DataInfoBucket, stock.ID)
	list := &InfoList{}
	if data == nil {
		return list
	}
	if err := xml.Unmarshal(data, list); err != nil {
		log.Println("Error: Failed to read data stock content", err)
		return &InfoList{}
	}
	return list
}

// UpdateContent updates the content of the given data stock.
func (db *DB) UpdateContent(stock *DataStock, content *InfoList) {
	data, err := xml.Marshal(content)
	if err != nil {
		log.Println("Error: Failed update data stock content", err)
		return
	}
	db.Put(DataInfoBucket, stock.ID, data)
}
