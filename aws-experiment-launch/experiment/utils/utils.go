package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
