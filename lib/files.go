package lib

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lnsp/filter"
)

const (
	DefaultFileMode = 0644
	PatchFile       = ".patch"
	IgnoreFile      = ".benchignore"
)

// Write creates a new file in the directory and puts a list of hashes in it.
// It may return an error if the file is not accessible.
func Write(dir, source string, data HashSet) error {
	target := filepath.Join(dir, PatchFile)
	joinBuffer, size := bytes.Buffer{}, len(data)
	if source != "" {
		joinBuffer.WriteString(SourceMarker + source + LineSeperator)
		log.Notice("generated patch with source", source)
	}
	for i := 0; i < size; i++ {
		joinBuffer.WriteString(data[i].String())
		joinBuffer.WriteString(LineSeperator)
	}
	outputFile, err := os.Create(target)
	if err != nil {
		log.Critical("failed to open patch file:", err)
		return err
	}
	defer outputFile.Close()
	joinBuffer.WriteTo(outputFile)
	return nil
}

// Generate creates a new patch file in the target folder storing version information.
// It may return an error if either the hashing or patch file creation fails.
func Generate(targetDir, targetSrc string) error {
	filterPath := filepath.Join(targetDir, IgnoreFile)
	filter := filter.LoadFilter(filterPath)
	hashes, err := HashDirectory(targetDir)
	if err != nil {
		log.Critical("failed to hash directory:", err)
		return err
	}
	filtered := FilterHashes(hashes, filter)
	return Write(targetDir, targetSrc, filtered)
}

// Fetch compares two patches from a global and a local branch and updates the local branch to match the global one.
// It only replaces or adds files, but does not delete any.
// It may return an error if the origin handling fails.
func Fetch(dir, target string) error {
	local, err := GetOrigin(dir)
	if err != nil {
		log.Error("bad local origin:", err)
		return err
	}
	localHashes, localSource, err := local.Scan()
	if target != "" {
		localSource = target
	}
	global, err := GetOrigin(localSource)
	if err != nil {
		log.Error("bad global origin:", err)
		return err
	}
	globalHashes, globalSource, err := global.Scan()
	if globalSource != localSource {
		log.Warning("unverified origin:", globalSource)
	} else {
		log.Notice("verified origin:", globalSource)
	}

	missingHashes := Compare(localHashes, globalHashes)
	for _, hash := range missingHashes {
		data, err := global.Get(hash.Name)
		if err != nil {
			log.Error("failed to fetch file:", hash.Name)
			return err
		}
		file := filepath.Join(dir, hash.Name)
		fileDir := filepath.Dir(file)
		err = os.MkdirAll(fileDir, DefaultFileMode)
		if err != nil {
			log.Error("failed to create folder:", fileDir)
			return err
		}
		ioutil.WriteFile(file, data, DefaultFileMode)
	}

	log.Notice("fetched", len(missingHashes), "files from origin")

	return Write(dir, globalSource, globalHashes)
}
