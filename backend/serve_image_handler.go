package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func newServeImageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("image")
		if imageName == "" {
			log.Printf("No image name given")
			return c.String(400, "No image name given")
		}
		responseWriter := c.Response().Writer
		responseWriter.Header().Set("Content-Type", "image/png")
		fileFullyQualified := filepath.Join("imagefiles", imageName)
		file, err := os.Open(fileFullyQualified)
		if err != nil {
			log.Printf("ERROR opening file %s", fileFullyQualified)
			return err
		}
		defer file.Close()

		gzReader, err := gzip.NewReader(file)
		if err != nil {
			log.Printf("ERROR with gzip decompressing: %s", err)
			return err
		}
		defer gzReader.Close()
		_, err = io.Copy(responseWriter, gzReader)
		if err != nil {
			log.Printf(
				"Error copying image (filename: %s) down to the response: %s",
				fileFullyQualified,
				err,
			)

			return err
		}
		return nil
	}
}
