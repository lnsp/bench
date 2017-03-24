// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Bench is a file patching system using HTTPS and secure file hashing.
package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/lnsp/bench/lib"
	"github.com/lnsp/pkginfo"
)

var (
	// DynamicPool controls the state of the worker pool allocation.
	DynamicPool = flag.Bool("dynamic", true, "Dynamic worker pool")
	// PoolSize sets the target worker count.
	PoolSize = flag.Int("worker", 1, "Async worker count")
	// PatchSource sets the source target.
	PatchSource = flag.String("source", "./", "Source target")
	// TargetDirectory sets the targeted folder.
	TargetDirectory = flag.String("target", "./", "Local target")
	// VerboseOutput enables debug logging output.
	VerboseOutput = flag.Bool("verbose", false, "Verbose logging")
	// PkgInfo stores information about the software package.
	PkgInfo = pkginfo.PackageInfo{
		Name: "bench",
		Version: pkginfo.PackageVersion{
			Major: 0,
			Minor: 2,
			Patch: 1,
		},
	}
)

const (
	// HelpText stores help text shown when running bench without parameters.
	HelpText = `USAGE: bench [action]

Available actions:
	generate - Generate patch files from active folder
	version - Print software version information
	fetch - Fetch updated files using file or server origin
	help - Display command overview

Available flags:
	--target - Local target (default "./")
	--verbose - Verbose logging (default false)
	--source - Source target (default "./")
	--worker - Async worker count (default 1)
	--dynamic - Dynamic worker count (default true)`
)

func main() {
	flag.Parse()

	lib.SetVerbose(*VerboseOutput)
	workingDir, _ := filepath.Abs(*TargetDirectory)

	command := flag.Arg(0)
	switch command {
	case "generate":
		lib.Generate(workingDir, *PatchSource, *PoolSize, *DynamicPool)
	case "fetch":
		lib.Fetch(workingDir, *PatchSource, *PoolSize, *DynamicPool)
	case "version":
		printVersion()
	case "help":
		fallthrough // to existing help case
	default:
		fmt.Println(HelpText)
	}
}

func printVersion() {
	fmt.Println(PkgInfo)
}
