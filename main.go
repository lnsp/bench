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
	PoolSize        = flag.Int("worker", 1, "Async worker count")
	PatchSource     = flag.String("source", "", "Set patch source")
	TargetDirectory = flag.String("target", ".", "Set custom working directory")
	VerboseOutput   = flag.Bool("verbose", false, "Enable verbose output")
	PkgInfo         = pkginfo.PackageInfo{
		Name: "bench",
		Version: pkginfo.PackageVersion{
			Major:      0,
			Minor:      2,
			Patch:      0,
			Identifier: "dev",
		},
	}
)

const (
	HelpText = `USAGE: bench [action]

Available actions:
	generate - Generate patch files from active folder
	version - Print software version information
	fetch - Fetch updated files using file or server origin
	help - Display command overview

Available flags:
	--target - Set custom working directory
	--verbose - Enable verbose output
	--source - Set patch source
	--worker - Async worker count (default 1)`
)

func main() {
	flag.Parse()

	lib.SetVerbose(*VerboseOutput)
	workingDir, _ := filepath.Abs(*TargetDirectory)

	command := flag.Arg(0)
	switch command {
	case "generate":
		lib.Generate(workingDir, *PatchSource)
	case "fetch":
		lib.Fetch(workingDir, *PatchSource, *PoolSize)
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
