package main

import (
	echo "github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
)

func newBasicStatsHandler(db *bolt.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		numImages := 0
		err := db.View(func(tx *bolt.Tx) error {
			nameLookupBucket := tx.Bucket([]byte(
				imageLookupBucketName,
			))

			return nameLookupBucket.ForEach(func(k, v []byte) error {
				numImages++
				return nil
			})
			return nil
		})
		if err != nil {
			c.Logger().Warnf("Error on DB transaction: %s", err)
			return err
		}
		return c.JSON(200, map[string]int{
			"num_images": numImages,
		})
	}
}
