package main

import (
	"log"

	echo "github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
)

const dbFile = "images.db"
const imageLookupBucketName = "image_name_lookup"

func main() {
	db, err := bolt.Open(dbFile, 0666, nil)
	if err != nil {
		log.Fatalf("Error opening database file: %s", err)
	}

	createBucketErr := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(
			imageLookupBucketName,
		))
		return err
	})
	if createBucketErr != nil {
		log.Fatal("Error creating image lookup bucket")
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, "This is the index page")
	})
	e.GET("/basic_stats", newBasicStatsHandler(db))

	e.POST("/process_image", newProcessHandler(db))

	e.GET("/image/:image", newImageHandler(db))
	e.Static("/serve_image", "./imagefiles")

	e.Logger.Fatal(e.Start(":5001"))
}
