// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lnsp/go-filter"
)

const (
	// HashSplit is the split symbol in each pair of name / hash sets.
	HashSplit = ":"
	// SourceMarker is the prefix of a source identifier in a patch file.
	SourceMarker = "#@"
	// LineSeparator is the default line Separator.
	LineSeparator = "\n"
)

// HashSet is a collection of hash items.
type HashSet []HashItem

// HashItem stores a hash with its identifier.
type HashItem struct {
	Name, Hash string
}

// String generates a valid string representation of the item.
func (hs HashItem) String() string {
	return hs.Name + HashSplit + hs.Hash
}

// Parse parses a patch file and returns the parsed items.
func Parse(data string) (HashSet, string) {
	items := strings.Split(data, LineSeparator)
	hashItems := make(HashSet, 0, len(items))
	source := ""
	for _, e := range items {
		if strings.HasPrefix(e, SourceMarker) {
			source = strings.TrimPrefix(e, SourceMarker)
			log.Info("found source in patch:", source)
		}
		tokens := strings.Split(strings.TrimSpace(e), HashSplit)
		// ignore invalid entry
		if len(tokens) < 2 {
			continue
		}
		hashItems = append(hashItems, HashItem{tokens[0], tokens[1]})
	}
	log.Info("parsed", len(hashItems), "hash items")
	return hashItems, source
}

// HashWorker waits for files, hashes them and sends back the result.
func HashWorker(base string, files <-chan string, results chan<- *HashItem) {
	for file := range files {
		hash, err := HashFile(filepath.Join(base, file))
		if err != nil {
			results <- nil
		} else {
			results <- &HashItem{
				Name: file,
				Hash: hash,
			}
		}
	}
}

// HashDirectoryAsync walks a directory recursively and generates a slice of hash pairs.
// A hash pair consists of the relative path from the start directory to the file
// and the hash of the file separated by a split string.
// It may return an error if the file hashing or directory walking fails.
func HashDirectoryAsync(start string, pool int) (HashSet, error) {
	var set HashSet
	jobs, err := ListFiles(start)
	if err != nil {
		return set, err
	}

	offset, workload := 0, len(jobs)
	set = make(HashSet, workload)
	files := make(chan string, workload/2+1)
	results := make(chan *HashItem, workload/2+1)

	for i := 0; i < pool; i++ {
		go HashWorker(start, files, results)
	}

	for i := 0; i < workload; i++ {
		files <- jobs[i]
	}
	close(files)

	for i := 0; i < workload; i++ {
		r := <-results
		if r == nil {
			offset++
			continue
		} else {
			set[i-offset] = *r
		}
	}

	return set[:(workload - offset)], nil
}

// HashDirectory walks a directory recursively and generates a slice of hash pairs.
// A hash pair consists of the relative path from the start directory to the file
// and the hash of the file separated by a split string.
// It may return an error if the file hashing or directory walking fails.
func HashDirectory(start string) (HashSet, error) {
	var set HashSet
	files, err := ListFiles(start)
	if err != nil {
		return set, err
	}

	offset, workload := 0, len(files)
	set = make(HashSet, workload)
	for i, element := range files {
		path := filepath.Join(start, element)
		hash, err := HashFile(path)
		if err != nil {
			offset++
			continue
		}
		set[i-offset] = HashItem{element, hash}
	}

	return set[:(workload - offset)], nil
}

// HashFile return a hex-encoded SHA-1 hash string of the file.
// It may return an error if the file is not found.
func HashFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Warning("file not found:", err)
		return "", err
	}

	hexString := hex.EncodeToString(GenerateHash(data))
	return hexString, nil
}

// GenerateHash generates a SHA-1 hash of the input data.
func GenerateHash(data []byte) []byte {
	sha := sha1.New()
	sha.Write(data)
	hashBytes := sha.Sum(nil)
	return hashBytes
}

// FilterHashes applies the given filter to the name of each element and
// returns a slice of the filtered items.
func FilterHashes(h HashSet, f filter.Filter) HashSet {
	filtered := make(HashSet, 0)
	for _, e := range h {
		if !f.Match(e.Name) {
			filtered = append(filtered, e)
		}
	}
	log.Info("ignored", len(h)-len(filtered), "files")
	return filtered
}

// Compare compares two slices of hash items and returns the difference between the two.
func Compare(local, global HashSet) HashSet {
	missingLocal := make(map[string]bool)
	log.Info("comparing branches: local [", len(local), "] <-> global [", len(global), "]")
	for _, e := range global {
		missingLocal[e.String()] = true
	}
	for _, e := range local {
		if _, ok := missingLocal[e.String()]; ok {
			missingLocal[e.String()] = false
		}
	}
	missingFiles := make(HashSet, 0)
	for _, e := range global {
		if missingLocal[e.String()] {
			missingFiles = append(missingFiles, e)
		}
	}
	log.Notice("missing in local branch:", len(missingFiles))

	return missingFiles
}
