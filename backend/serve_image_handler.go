package main

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func newServeImageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		imageName := c.Param("image")
		if imageName == "" {
			return c.String(400, "No image name given")
		}
		responseWriter := c.Response().Writer
		responseWriter.Header().Set("Content-Type", "image/x-icon")
		fileFullyQualified := filepath.Join("imagefiles", imageName)
		file, err := os.Open(fileFullyQualified)
		if err != nil {
			return err
		}
		defer file.Close()

		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(responseWriter, gzReader)
		if err != nil {
			return err
		}
		return nil
	}
}
