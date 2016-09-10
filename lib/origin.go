// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package lib

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// Origin is a accessable data store.
type Origin interface {
	// Scan fetches the stored HashSet and source name from the origin.
	// It may return an error if the set is not accessible.
	Scan() (HashSet, string, error)
	// Get fetches a file specified by its path from the origin.
	// It may return an error if the file is not accessible.
	Get(file string) ([]byte, error)
}

// HTTPOrigin has a HTTP back-end.
type HTTPOrigin struct {
	Base string
}

func (h HTTPOrigin) Get(file string) ([]byte, error) {
	requestURL, err := url.Parse(h.Base)
	if err != nil {
		log.Error("url parse:", err)
		return []byte{}, err
	}
	requestURL.Path = filepath.Join(requestURL.Path, file)

	resp, err := http.Get(requestURL.String())
	if err != nil {
		log.Warning("http request:", err)
		return []byte{}, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warning("http transmission:", err)
		return []byte{}, err
	}
	return buf, nil
}

func (h HTTPOrigin) Scan() (HashSet, string, error) {
	data, err := h.Get(PatchFile)
	if err != nil {
		log.Warning("http scan:", err)
		return HashSet{}, "", err
	}
	items, source := Parse(string(data))
	log.Notice("fetched", len(items), "from http origin")
	return items, source, nil
}

// FileOrigin has a filesystem back-end.
type FileOrigin struct {
	Path string
}

func (origin FileOrigin) Scan() (HashSet, string, error) {
	data, err := origin.Get(PatchFile)
	if err != nil {
		log.Warning("local scan:", err)
		return HashSet{}, "", nil
	}
	items, source := Parse(string(data))
	log.Notice("fetched", len(items), "from local origin")
	return items, source, nil
}

func (origin FileOrigin) Get(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(origin.Path, file))
	if err != nil {
		log.Warning("local file:", err)
		return []byte{}, err
	}
	return data, nil
}

// GetOrigin finds a fitting origin based on the path.
// It may return an error if there if the path pattern is unknown.
func GetOrigin(path string) (Origin, error) {
	if strings.HasPrefix(path, "http") {
		return &HTTPOrigin{path}, nil
	} else if strings.HasPrefix(path, "/") {
		return &FileOrigin{path}, nil
	} else {
		return nil, errors.New("unknown origin: " + path)
	}
}
