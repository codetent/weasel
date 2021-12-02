package store

import (
	"fmt"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const (
	BucketDistributions string = "distributions"
	BucketInstances     string = "instances"
)

type RegisteredDistribution struct {
	Id string
}

type RegisteredInstance struct {
	Id           string
	Distribution string
}

func OpenDatabase() (*bolt.DB, error) {
	cacheRoot, err := GetCacheRoot()
	if err != nil {
		return nil, fmt.Errorf("OpenDatabase: GetCacheRoot(): %v", err)
	}

	dbPath := filepath.Join(cacheRoot, "weasel.db")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("OpenDatabase: Open(): %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{BucketDistributions, BucketInstances} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket '%s': %v", bucket, err)
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("OpenDatabase: Update(): %v", err)
	}

	return db, nil
}

func RegisterDistribution(id string, path string) error {
	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("RegisterDistribution: OpenDatabase(): %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketDistributions))
		return bucket.Put([]byte(id), []byte(path))
	})
	if err != nil {
		return fmt.Errorf("RegisterDistribution: Update(): %v", err)
	}

	return nil
}

func UnregisterDistribution(id string) error {
	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("UnregisterDistribution: OpenDatabase(): %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketDistributions))
		return bucket.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("UnregisterDistribution: Update(): %v", err)
	}

	return nil
}

func GetRegisteredDistribution(id string) (string, error) {
	db, err := OpenDatabase()
	if err != nil {
		return "", fmt.Errorf("GetRegisteredDistribution: OpenDatabase(): %v", err)
	}
	defer db.Close()

	rx, err := db.Begin(false)
	if err != nil {
		return "", fmt.Errorf("GetRegisteredDistribution: Begin(): %v", err)
	}
	defer rx.Rollback()

	bucket := rx.Bucket([]byte(BucketDistributions))
	raw := bucket.Get([]byte(id))
	return string(raw), nil
}

func GetRegisteredDistributions() ([]RegisteredDistribution, error) {
	db, err := OpenDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetRegisteredDistribution: OpenDatabase(): %v", err)
	}
	defer db.Close()

	rx, err := db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("GetRegisteredDistributions: Begin(): %v", err)
	}
	defer rx.Rollback()

	bucket := rx.Bucket([]byte(BucketDistributions))
	cursor := bucket.Cursor()
	var dists []RegisteredDistribution

	for id, _ := cursor.First(); id != nil; id, _ = cursor.Next() {
		dists = append(dists, RegisteredDistribution{string(id)})
	}

	return dists, nil
}

func RegisterInstance(id string, dist string) error {
	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("RegisterInstance: OpenDatabase(): %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketInstances))
		return bucket.Put([]byte(id), []byte(dist))
	})
	if err != nil {
		return fmt.Errorf("RegisterInstance: Update(): %v", err)
	}

	return nil
}

func UnregisterInstance(id string) error {
	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("UnregisterInstance: OpenDatabase(): %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketInstances))
		return bucket.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("UnregisterInstance: Update(): %v", err)
	}

	return nil
}

func GetRegisteredInstances() ([]RegisteredInstance, error) {
	db, err := OpenDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetInstanceIds: OpenDatabase(): %v", err)
	}
	defer db.Close()

	rx, err := db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("GetInstanceIds: Begin(): %v", err)
	}
	defer rx.Rollback()

	bucket := rx.Bucket([]byte(BucketInstances))
	cursor := bucket.Cursor()
	var instances []RegisteredInstance

	for id, dist := cursor.First(); id != nil; id, dist = cursor.Next() {
		instances = append(instances, RegisteredInstance{string(id), string(dist)})
	}

	return instances, nil
}
