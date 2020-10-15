package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
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

	e.POST("/process_image", newProcessHandler(db))
}

func newProcessHandler(db *bolt.DB) echo.HandlerFunc {
	type req struct {
		URL  string   `json:"url"`
		Tags []string `json:"tags"`
		Name string   `json:"name"`
	}
	return func(c echo.Context) error {
		r := &req{}
		if err := c.Bind(r); err != nil {
			return err
		}

		// Go's DB support lets us wrap our entire HTTP handler
		// inside of a transaction. This is similar to a SQLIte
		// transaction if you're using Python
		db.Update(func(tx *bolt.Tx) error {
			// first, check if the image name exists. that would
			// be indicated if the bucket for that image
			// currently exists
			imgFilename := fmt.Sprintf("%s-%s", r.Name, uuid.New())
			imageBucketName := []byte(imgFilename)

			if tx.Bucket(imageBucketName) != nil {
				return fmt.Errorf("Image name is taken")
			}

			// next, download and store the image to a file
			imgFile, err := os.OpenFile(imgFilename, os.O_WRONLY, 0755)
			if err != nil {
				return err
			}

			// Go lets us specify a different HTTP client for specific
			// cases. In this case, the default one is fine, but
			// we have the option of gaining very low level control
			// over connections, timeouts, etc... if we want
			// to create our own client.
			//
			// Create a request to fetch the image, and then download
			// it
			cl := http.DefaultClient
			req, err := http.NewRequest("GET", r.URL, nil)
			if err != nil {
				return err
			}
			res, err := cl.Do(req)
			if err != nil {
				return err
			}
			if res.StatusCode >= 400 {
				return fmt.Errorf(
					"Couldn't get image from %s, status code %d",
					r.URL,
					res.StatusCode,
				)
			}

			defer res.Body.Close()

			// Next, store the downloaded file onto disk.
			// Go has the ability to "compose" I/O operations together
			// through it's standard library
			// e can stream
			// the data from the origin server that the image is hosted on,
			// through the compressor, and directly down to the file.
			// Depending on the
			// down to the file, reducing the amount of memory our
			// program needs to use.
			//
			// this is an example of how Go allows you to take control of
			// more of the low-level features of the underlying operating system
			// in order to optimize the performance of your program
			gzipWriter := gzip.NewWriter(imgFile)

			// here's we're streaming data from the response body
			// (which may be streaming directly from the server)
			// through the Gzip compressor, then down to the actual
			// on-disk file
			imageSize, err := io.Copy(gzipWriter, res.Body)
			if err != nil {
				return err
			}

			// next, write all of the file metadata to the database.
			// remember, we're doing all of this inside a transaction, so
			// if these operations or anything before fails, we will
			// be rolling back all of the operations we've done inside
			// here

			// laying out the DB like this:
			// - one bucket per image.
			//		- bucket name is the same uuid as what is in the image filename
			//		- keys: size, name, tags, filename, url
			// - one bucket for the name -> image bucket lookup table

			imageBucket, err := tx.CreateBucketIfNotExists(imageBucketName)
			if err != nil {
				return err
			}

			imageMetadata := map[string]interface{}{
				"size":     imageSize,
				"name":     r.Name,
				"tags":     r.Tags,
				"filename": imgFilename,
				"url":      r.URL,
			}
			for k, v := range imageMetadata {
				valBytes, err := json.Marshal(v)
				if err != nil {
					return err
				}
				imageBucket.Put([]byte(k), valBytes)
			}

			nameLookupBucket, err := tx.CreateBucketIfNotExists([]byte(
				imageLookupBucketName,
			))
			if err != nil {
				return err
			}

			if err := nameLookupBucket.Put([]byte(r.Name), []byte(
				imageBucketName,
			)); err != nil {
				return err
			}

			return nil

		})
		return nil
	}
}
