// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package lib generates and fetches hash patches.
package lib

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lnsp/filter"
)

const (
	// DefaultFileMode is the default file and directory mode.
	DefaultFileMode = 0644
	// PatchFile is the default patch file name.
	PatchFile = ".patch"
	// IgnoreFile is the default ignore file name.
	IgnoreFile = ".benchignore"
)

// Write creates a new file in the directory and puts a list of hashes in it.
// It may return an error if the file is not accessible.
func Write(dir, source string, data HashSet) error {
	target := filepath.Join(dir, PatchFile)
	joinBuffer, size := bytes.Buffer{}, len(data)
	// append source if it exists
	if source != "" {
		joinBuffer.WriteString(SourceMarker + source + LineSeperator)
		log.Notice("generated patch with source", source)
	}
	// store all datasets in the file
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
	// copy byte buffer into file
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

// FetchSpecific downloads the files in the set from the origin and stores it in the target directory.
// It may fail if either can't fetch the file or can't create the directory / file.
func FetchSpecific(dir string, source Origin, set HashSet) error {
	for _, hash := range set {
		data, err := source.Get(hash.Name)
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
		err = ioutil.WriteFile(file, data, DefaultFileMode)
		if err != nil {
			log.Error("failed to write file:", file)
			return err
		}
	}
	return nil
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
	err = FetchSpecific(dir, global, missingHashes)
	if err != nil {
		log.Error("failed to fetch files:", err)
		return err
	}
	log.Notice("fetched", len(missingHashes), "files from origin")

	return Write(dir, globalSource, globalHashes)
}
