package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/process_image", newProcessHandler())
}

func newProcessHandler() echo.HandlerFunc {
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

		// 1. download and store image
		imgFilename := fmt.Sprintf("%s-%s", r.Name, uuid.New())
		imgFile, err := os.OpenFile(imgFilename, os.O_WRONLY, 0755)
		if err != nil {
			return err
		}

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
		// we can stream the file from the remote server
		// onto disk here, reducing the amount of memory our
		// program needs to use. this is an example of
		// how Go allows you to take control of more of the
		// low-level features of the underlying operating system
		// in order to optimize the performance of your program
		if _, err := io.Copy(imgFile, res.Body); err != nil {
			return err
		}

		// 2. store tags, name, and file location to a DB somewhere
		//		strong support for structured DBs, written in native
		//		Go, no extra things to compile

		// we're using boltdb for the image metadata database
		// and tags index
		//
		// laying out the DB like this:
		// - one bucket per image.
		//		- bucket name is the same uuid as what is in the image filename
		//		- keys: size, name, tags, filename, url
		// - one bucket for global stats:
		//		- keys: num images, total size, average size, avg tags per image
		// - one bucket for the tag index
		//		- each key is a tag, value is a list of image bucket keys that have that tag

		return nil
	}
}
