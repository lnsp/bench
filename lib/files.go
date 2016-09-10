// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package lib generates and fetches hash patches.
package lib

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

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
func Generate(targetDir, targetSrc string, pool int, dynamic bool) error {
	filterPath := filepath.Join(targetDir, IgnoreFile)
	filter := filter.LoadFilter(filterPath)
	var hashes HashSet
	var err error
	if pool < 2 {
		hashes, err = HashDirectory(targetDir)
	} else {
		if dynamic {
			pool *= runtime.NumCPU()
		}
		hashes, err = HashDirectoryAsync(targetDir, pool)
	}
	if err != nil {
		log.Critical("failed to hash directory:", err)
		return err
	}
	filtered := FilterHashes(hashes, filter)
	return Write(targetDir, targetSrc, filtered)
}

// FetchWorker is an async worker waiting for jobs.
func FetchWorker(elements <-chan HashItem, results chan<- error, origin Origin, dir string) {
	for e := range elements {
		data, err := origin.Get(e.Name)
		if err != nil {
			results <- errors.New("failed to fetch file: " + e.Name)
			continue
		}
		file := filepath.Join(dir, e.Name)
		fileDir := filepath.Dir(file)
		err = os.MkdirAll(fileDir, DefaultFileMode)
		if err != nil {
			results <- errors.New("failed to create folder: " + filepath.Dir(e.Name))
			continue
		}
		err = ioutil.WriteFile(file, data, DefaultFileMode)
		if err != nil {
			results <- errors.New("failed to write file: " + e.Name)
			continue
		}
		results <- nil
	}
}

// FetchSpecificAsync downloads the files asynchronously in the set from the origin and stores it in the target directory.
// It may fail if either can't fetch the file or can't create the directory / file.
func FetchSpecificAsync(target string, source Origin, set HashSet, pool int) error {
	workload := len(set)
	jobs := make(chan HashItem, workload/2+1)
	results := make(chan error, workload/2+1)

	log.Notice("using", pool, "workers")

	for i := 0; i < pool; i++ {
		go FetchWorker(jobs, results, source, target)
	}

	for i := 0; i < workload; i++ {
		jobs <- set[i]
	}
	close(jobs)

	var err error
	for i := 0; i < workload; i++ {
		if err = <-results; err != nil {
			log.Warning(err)
		}
	}

	return nil
}

// FetchSpecific downloads the files in the set from the origin and stores it in the target directory.
// It may fail if either can't fetch the file or can't create the directory / file.
func FetchSpecific(dir string, source Origin, set HashSet) error {
	for _, hash := range set {
		data, err := source.Get(hash.Name)
		if err != nil {
			return errors.New("failed to fetch file: " + hash.Name)
		}
		file := filepath.Join(dir, hash.Name)
		fileDir := filepath.Dir(file)
		err = os.MkdirAll(fileDir, DefaultFileMode)
		if err != nil {
			return errors.New("failed to create folder: " + fileDir)
		}
		err = ioutil.WriteFile(file, data, DefaultFileMode)
		if err != nil {
			return errors.New("failed to write file: " + file)
		}
	}
	return nil
}

// Fetch compares two patches from a global and a local branch and updates the local branch to match the global one.
// It only replaces or adds files, but does not delete any.
// It may return an error if the origin handling fails.
func Fetch(dir, target string, pool int, dynamic bool) error {
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
	if pool < 2 {
		err = FetchSpecific(dir, global, missingHashes)
	} else {
		if dynamic {
			pool *= runtime.NumCPU()
		}
		err = FetchSpecificAsync(dir, global, missingHashes, pool)
	}
	if err != nil {
		log.Error("failed to fetch files:", err)
		return err
	}
	log.Notice("fetched", len(missingHashes), "files from origin")

	return Write(dir, globalSource, globalHashes)
}
