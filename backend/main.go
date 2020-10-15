package main

import (
	"log"

	"github.com/boltdb/bolt"
	echo "github.com/labstack/echo/v4"
)

const dbFile = "images.db"
const imageLookupBucketName = "image_name_lookup"

func main() {
	db, err := bolt.Open(dbFile, 0666, nil)
	if err != nil {
		log.Fatalf("Error opening database file: %s", err)
	}

	e := echo.New()

	e.GET("/basic_stats", newBasicStatsHandler(db))

	e.POST("/process_image", newProcessHandler(db))
	e.Logger.Fatal(e.Start(":5001"))
}
