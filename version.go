// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"fmt"
	"regexp"
	"strconv"
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
//
// It follows semver specification (e.g. v1.200.300 is less than v2). A dev
// version is *always* less than a non-dev version (e.g. v3-dev is less than
// v2).
func (v Version) Less(other Version) bool {
	if v.Dev && !other.Dev {
		return true
	} else if other.Dev && !v.Dev {
		return false
	}

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

// Matches strings like "1", "1.1", and "1.1.1".
var vsRegexp = regexp.MustCompile(`^([0-9]+)[\.]?([0-9]*)[\.]?([0-9]*)`)

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

	// Split by the dash seperated suffix. We expect only one dash suffix, and
	// if present it must be "dev".
	dashSplit := strings.Split(vs, "-")
	if len(dashSplit) > 2 || len(dashSplit) == 2 && dashSplit[1] != "dev" {
		return InvalidVersion
	}

	// We now use regexp to match the last part of the version string, which
	// e.g. looks like "1", "1.1", or "1.1.1".
	var (
		m = vsRegexp.FindStringSubmatch(dashSplit[0])
		v = InvalidVersion
	)
	if len(m) > 1 && len(m[1]) > 0 {
		v.Major, _ = strconv.Atoi(m[1])
	}
	if len(m) > 2 && len(m[2]) > 0 {
		v.Minor, _ = strconv.Atoi(m[2])
	}
	if len(m) > 3 && len(m[3]) > 0 {
		v.Patch, _ = strconv.Atoi(m[3])
	}
	if v.Major != -1 && len(dashSplit) == 2 && dashSplit[1] == "dev" {
		v.Dev = true
	}
	return v
}
