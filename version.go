// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"fmt"
	"strings"
)

// Version represents a semantic version.
type Version struct {
	Major, Minor, Patch int

	// If true, then this is the in-development version.
	Dev bool
}

// String returns a string representation of this version, for example:
//
//  Version{Major=1, Minor=2, Patch=3}           -> "v1.2.3"
//  Version{Major=1, Minor=2, Patch=3, Dev=true} -> "v1.2.3-dev"
//
//  Version{Major=1, Minor=2, Patch=-1}           -> "v1.2"
//  Version{Major=1, Minor=2, Patch=-1, Dev=true} -> "v1.2-dev"
//
//  Version{Major=1, Minor=-1, Patch=-1}           -> "v1"
//  Version{Major=1, Minor=-1, Patch=-1, Dev=true} -> "v1-dev"
//
func (v Version) String() string {
	var s string
	if v.Major > 0 && v.Minor > 0 && v.Patch > 0 {
		s = fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else if v.Major > 0 && v.Minor > 0 {
		s = fmt.Sprintf("v%d.%d", v.Major, v.Minor)
	} else if v.Major > 0 {
		s = fmt.Sprintf("v%d", v.Major)
	} else {
		return fmt.Sprintf("Version{Major=%d, Minor=%d, Patch=%d, Dev=%t}", v.Major, v.Minor, v.Patch, v.Dev)
	}
	if v.Dev {
		return s + "-dev"
	}
	return s
}

// Less tells if v is a lesser version than the other version.
func (v Version) Less(other Version) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major > other.Major {
		return false
	}

	if v.Minor < other.Minor {
		return true
	} else if v.Minor > other.Minor {
		return false
	}

	if v.Patch < other.Patch {
		return true
	} else if v.Patch > other.Patch {
		return false
	}
	return false
}

// InvalidVersion represents a completely invalid version.
var InvalidVersion = Version{
	Major: -1,
	Minor: -1,
	Patch: -1,
	Dev:   false,
}

// ParseVersion parses a version string in the form of:
//
//  "v1"
//  "v1.2"
//  "v1.2.1"
//  "v1-dev"
//  "v1.2-dev"
//  "v1.2.1-dev"
//
// It returns InvalidVersion for strings not suffixed with "v", like:
//
//  "1"
//  "1.2-dev"
//
func ParseVersion(vs string) Version {
	if vs[0] != 'v' {
		return InvalidVersion
	}
	vs = vs[1:] // Strip prefixed v

	// Parse -dev suffix.
	var v Version
	v.Dev = strings.HasSuffix(vs, "-dev")
	if v.Dev {
		vs = strings.TrimSuffix(vs, "-dev")
	}

	// Parse actual version number.
	switch strings.Count(vs, ".") {
	default:
		return InvalidVersion
	case 0:
		fmt.Sscanf(vs, "%d", &v.Major)
	case 1:
		fmt.Sscanf(vs, "%d.%d", &v.Major, &v.Minor)
	case 2:
		fmt.Sscanf(vs, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	}
	return v
}
