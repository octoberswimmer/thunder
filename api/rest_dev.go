//go:build dev
// +build dev

package api

import (
	"io"
	"net/http"
)

// Get performs an HTTP GET against the local dev server and returns the response body.
func Get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data
}
