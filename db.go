package main

import (
	"path/filepath"

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
func (db *DB) Get(id string, bucket BucketType) ([]byte, error) {
	var data []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data = b.Get([]byte(id))
		return nil
	})
	return data, err
}

// Delete deletes the entity with the given ID from the bucket.
func (db *DB) Delete(id string, bucket BucketType) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		return bucket.Delete([]byte(id))
	})
}
