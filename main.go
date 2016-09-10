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
	PatchSource     = flag.String("source", "", "Set patch source")
	TargetDirectory = flag.String("target", ".", "Set custom working directory")
	VerboseOutput   = flag.Bool("verbose", false, "Enable verbose output")
	PkgInfo         = pkginfo.PackageInfo{
		Name: "bench",
		Version: pkginfo.PackageVersion{
			Major:      0,
			Minor:      1,
			Patch:      0,
			Identifier: "release",
		},
	}
)

const (
	PatchFilePerm = 0664
	PatchFile     = ".bpatch"
	IgnoreFile    = ".benchignore"
	HelpText      = `USAGE: bench [action]

Available actions:
	generate - Generate patch files from active folder
	version - Print software version information
	fetch - Fetch updated files using file or server origin
	help - Display command overview

Available flags:
	--target - Set custom working directory
	--verbose - Enable verbose output
	--source - Set patch source`
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
		lib.Fetch(workingDir, *PatchSource)
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
