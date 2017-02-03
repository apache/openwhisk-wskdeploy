package utils

import (
	"io/ioutil"
	"net/http"
)

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}

type URLReader struct {
}

func (urlReader *URLReader) ReadUrl(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	Check(err)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Check(err)
	return b, nil
}

type LocalReader struct {
}

func (localReader *LocalReader) ReadLocal(path string) (content []byte, err error) {
	cont, err := ioutil.ReadFile(path)
	Check(err)
	return cont, nil
}
