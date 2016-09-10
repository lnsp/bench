// Package filter provides file path filtering.
package filter

import (
	"github.com/lnsp/filematch"
	"io/ioutil"
	"strings"
)

// Filter is a collection of rules.
type Filter []string

// Match tests if the string does fulfill any of the rules.
func (f Filter) Match(e string) bool {
	match, err := filematch.MatchPath(e, []string(f))
	if err != nil {
		return false
	}
	return match
}

// Add a rule to the collection.
func (f Filter) Add(part string) {
	f = append(f, part)
}

// LoadFilter creates a filter from the file.
func LoadFilter(file string) Filter {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		data = make([]byte, 0)
	}
	var entries []string
	for _, entry := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(entry)
		if strings.HasPrefix(trimmed, "#") || trimmed == "" {
			continue
		}
		entries = append(entries, trimmed)
	}
	return Filter(entries)
}
