// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"net/url"
	"testing"
)

var gitHubRegexpTests = []struct {
	full, name, version string
	valid               bool
}{
	{"pkg.v3.2.1", "pkg", "v3.2.1", true},
	{"pkg.v3.2", "pkg", "v3.2", true},
	{"pkg.v3", "pkg", "v3", true},
	{"pkg.v0", "pkg", "v0", true},
	{"go-pkg.v1", "go-pkg", "v1", true},
	{"etc.v1-unstable", "etc", "v1-unstable", true},
	{"etc.v3.2-unstable", "etc", "v3.2-unstable", true},
	{"abc/def.v3.2", "def", "v3.2", false},
	{"abc/other.v3/def", "other", "v3", false},
	{"pkg", "", "", false},
}

// Test the rePkgVersion regexp.
func TestGitHubRegexp(t *testing.T) {
	for _, tst := range gitHubRegexpTests {
		m := rePkgVersion.FindStringSubmatch(tst.full)

		// Validate the return value.
		if !tst.valid && m != nil {
			t.Logf("%q\n", tst.full)
			t.Fatal("want no match, got", m)
			continue
		} else if tst.valid && m == nil {
			t.Logf("%q\n", tst.full)
			t.Fatal("want match, got nil")
			continue
		}
		if !tst.valid {
			continue
		}

		// Validate length of matched data.
		if len(m) != 3 {
			t.Logf("%q\n", tst.full)
			for i, m := range m {
				t.Logf("%d. %q\n", i, m)
			}
			t.Fatal("expected 3 values, but regex matched", len(m))
			continue
		}

		// Validate the actual data.
		if m[1] != tst.name {
			t.Logf("%q\n", tst.full)
			t.Fatalf("expected name %q, got %q", tst.name, m[1])
			continue
		} else if m[2] != tst.version {
			t.Logf("%q\n", tst.full)
			t.Fatalf("expected version %q, got %q", tst.version, m[2])
			continue
		}
	}
}

var gitHubTests = []struct {
	url, github, subpath string
	valid                bool
}{
	{"pkg.v3", "bob/pkg.git", "", true},
	{"/pkg.v3", "bob/pkg.git", "", true},
	{"go-pkg.v4", "bob/go-pkg.git", "", true},
	{"folder/pkg.v3", "bob/folder-pkg.git", "", true},
	{"multi/folder/pkg.v3", "bob/multi-folder-pkg.git", "", true},
	{"folder/pkg.v3/subpkg", "bob/folder-pkg.git", "subpkg", true},
	{"pkg.v3/folder/subpkg", "bob/pkg.git", "folder/subpkg", true},
	{"go-pkg.v3/folder/subpkg", "bob/go-pkg.git", "folder/subpkg", true},
	{"a", "", "", false},
	{"a/b", "", "", false},
	{"a/b/", "", "", false},
	{"a/b.v3/", "", "", false},
	{"a.v3/b/c.v3", "", "", false},
}

// Tests the GitHub URL matcher.
func TestGitHub(t *testing.T) {
	for _, tst := range gitHubTests {
		// Parse test case URL.
		u, err := url.Parse(tst.url)
		if err != nil {
			t.Fatal(err)
		}

		// Create a matcher and perform matching.
		matcher := GitHub("bob")
		repo, err := matcher.Match(u)
		if tst.valid && err != nil {
			t.Log(u)
			t.Fatal("Test is valid but matcher returned:", err)
		} else if !tst.valid && err == nil {
			t.Log(u)
			t.Fatal("Test is invalid but matcher returned nil error!")
		}
		if !tst.valid {
			continue
		}

		// Validate repo structure.
		if repo == nil {
			t.Log(u)
			t.Fatal("nil repo")
		}
		if repo.URL.Path != tst.github {
			t.Log(u)
			t.Log("want", tst.github)
			t.Log("got", repo.URL.Path)
			t.Fatal("incorrect path")
		}
		if repo.SubPath != tst.subpath {
			t.Log(u)
			t.Log("want", tst.subpath)
			t.Log("got", repo.SubPath)
			t.Fatal("incorrect path")
		}
	}
}
