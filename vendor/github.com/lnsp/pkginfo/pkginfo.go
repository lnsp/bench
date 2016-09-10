// Copyright 2016 Lennart Espe. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package pkginfo provides utilities to store package metadata.
package pkginfo

import "fmt"

// PackageInfo stores meta information about a package like name and version number.
type PackageInfo struct {
	Name    string
	Version PackageVersion
}

// PackageVersion stores software versions using the semver (see semver.org for more information) standard.
type PackageVersion struct {
	Major, Minor, Patch int
	Identifier          string
	Build               string
}

func (pkgVersion PackageVersion) String() string {
	base := fmt.Sprintf("%d.%d.%d", pkgVersion.Major, pkgVersion.Minor, pkgVersion.Patch)
	if pkgVersion.Identifier != "" {
		base += "-" + pkgVersion.Identifier
	}
	if pkgVersion.Build != "" {
		base += "+" + pkgVersion.Build
	}
	return base
}

func (pkgInfo PackageInfo) String() string {
	return fmt.Sprintf("%s %s", pkgInfo.Name, pkgInfo.Version.String())
}
