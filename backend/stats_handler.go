package main

import (
	"github.com/boltdb/bolt"
	echo "github.com/labstack/echo/v4"
)

func newBasicStatsHandler(db *bolt.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		numImages := 0
		err := db.View(func(tx *bolt.Tx) error {
			nameLookupBucket, err := tx.CreateBucketIfNotExists([]byte(
				imageLookupBucketName,
			))
			if err != nil {
				return err
			}

			return nameLookupBucket.ForEach(func(k, v []byte) error {
				numImages++
				return nil
			})
			return nil
		})
		if err != nil {
			return err
		}
		return c.JSON(200, map[string]int{
			"num_images": numImages,
		})
	}
}
