package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lnsp/filter"
)

const (
	HashSplit     = ":"
	SourceMarker  = "#@"
	LineSeperator = "\n"
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
	items := strings.Split(data, LineSeperator)
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

// HashDirectory walks a directory recursively and generates a slice of hash pairs.
// A hash pair consists of the relative path from the start directory to the file
// and the hash of the file seperated by a split string.
// It may return an error if the file hashing or directory walking fails.
func HashDirectory(start string) (HashSet, error) {
	var hashPairs HashSet
	err := filepath.Walk(start, func(active string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(start, active)
		if err != nil {
			log.Warning(err)
			return err
		}
		hash, err := HashFile(active)
		if err != nil {
			log.Warning(err)
			return err
		}
		hashPairs = append(hashPairs, HashItem{relPath, hash})
		return nil
	})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return hashPairs, err
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

// generateHash generates a SHA-1 hash of the input data.
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
