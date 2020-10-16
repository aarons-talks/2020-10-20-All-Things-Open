package main

import (
	"fmt"
	"strings"

	echo "github.com/labstack/echo/v4"
	bolt "go.etcd.io/bbolt"
)

func newImageHandler(db *bolt.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("image")
		if imageName == "" {
			return c.JSON(
				400,
				fmt.Sprintf("No name given"),
			)
		}
		// first look in the global image bucket lookup table to get
		// the specific bucket name for the image
		imageBucketName := []byte{}
		err := db.View(func(tx *bolt.Tx) error {
			nameLookupBucket := tx.Bucket([]byte(
				imageLookupBucketName,
			))
			if nameLookupBucket == nil {
				return fmt.Errorf("Image %s not found", imageName)
			}
			imageBucketName = nameLookupBucket.Get([]byte(imageName))
			return nil
		})
		if err != nil {
			return c.JSON(
				404,
				fmt.Sprintf("Image %s not found (%s)", imageName, err),
			)
		}

		// next look up the specific image bucket to get the
		// filename of the image itself.
		filename := []byte{}
		err = db.View(func(tx *bolt.Tx) error {
			imageBucket := tx.Bucket(imageBucketName)
			if imageBucket == nil {
				return fmt.Errorf("Image %s not found", imageName)
			}

			filename = imageBucket.Get([]byte("filename"))
			if filename == nil {
				return fmt.Errorf("Filename was not available in the image bucket")
			}
			return nil
		})

		if err != nil {
			return c.JSON(
				500,
				fmt.Sprintf("Error lookup up the image bucket: %s", err),
			)
		}

		cleanFilename := strings.Trim(string(filename), `"`)
		// then, return some alt text for the image, along with
		// a URL back to this server (maybe a static file server??)
		// that the flask app can point the HTML to
		return c.JSON(418, map[string]string{
			"alt": imageName,
			"src": fmt.Sprintf("/serve_image/%s", cleanFilename),
		})
	}
}
