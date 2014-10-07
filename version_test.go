// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import "testing"

var versionLessTests = []struct {
	a, b  Version
	aLess bool
}{
	{
		a:     Version{Major: 1, Minor: 1, Patch: 1},
		b:     Version{Major: 1, Minor: 1, Patch: 1},
		aLess: false,
	}, {
		a:     Version{Major: 2, Minor: 1, Patch: 1},
		b:     Version{Major: 1, Minor: 1, Patch: 1},
		aLess: false,
	}, {
		a:     Version{Major: 0, Minor: 1, Patch: 1},
		b:     Version{Major: 1, Minor: 1, Patch: 1},
		aLess: true,
	}, {
		a:     Version{Major: 1, Minor: 1, Patch: 1},
		b:     Version{Major: 1, Minor: 2, Patch: 1},
		aLess: true,
	},
}

func TestVersionLess(t *testing.T) {
	for _, tst := range versionLessTests {
		less := tst.a.Less(tst.b)
		if less != tst.aLess {
			t.Log("a", tst.a)
			t.Log("b", tst.b)
			t.Fatal("a < b; got", less, "want", tst.aLess)
		}
	}
}

var versionParseTests = []struct {
	v   string
	exp Version
}{
	// Valid version strings.
	{v: "v1", exp: Version{Major: 1, Minor: -1, Patch: -1}},
	{v: "v1.2", exp: Version{Major: 1, Minor: 2, Patch: -1}},
	{v: "v1.2.3", exp: Version{Major: 1, Minor: 2, Patch: 3}},
	{v: "v100-dev", exp: Version{Major: 100, Minor: -1, Patch: -1, Dev: true}},
	{v: "v1.24-dev", exp: Version{Major: 1, Minor: 24, Patch: -1, Dev: true}},
	{v: "v14.2.34-dev", exp: Version{Major: 14, Minor: 2, Patch: 34, Dev: true}},

	// Version strings must have 'v' prefix.
	{v: "1", exp: InvalidVersion},
	{v: "1.2", exp: InvalidVersion},
	{v: "1.2.3", exp: InvalidVersion},
	{v: "100-dev", exp: InvalidVersion},
	{v: "1.24-dev", exp: InvalidVersion},
	{v: "14.2.34-dev", exp: InvalidVersion},

	// Invalid version strings.
	{v: "v-dev", exp: InvalidVersion},
	{v: "ga.v1.r.3.ba.4.ge", exp: InvalidVersion},
}

func TestVersionParsing(t *testing.T) {
	for _, tst := range versionParseTests {
		got := ParseVersion(tst.v)
		want := tst.exp
		if got != want {
			t.Logf("%q\n", tst.v)
			t.Logf("got Major=%d Minor=%d Patch=%d Dev=%t\n", got.Major, got.Minor, got.Patch, got.Dev)
			t.Fatalf("want Major=%d Minor=%d Patch=%d Dev=%t\n", want.Major, want.Minor, want.Patch, want.Dev)
		}
	}
}

var versionStringTests = []struct {
	v   string
	exp Version
}{
	// Valid version strings.
	{v: "v1", exp: Version{Major: 1, Minor: -1, Patch: -1}},
	{v: "v1.2", exp: Version{Major: 1, Minor: 2, Patch: -1}},
	{v: "v1.2.3", exp: Version{Major: 1, Minor: 2, Patch: 3}},
	{v: "v100-dev", exp: Version{Major: 100, Minor: -1, Patch: -1, Dev: true}},
	{v: "v1.24-dev", exp: Version{Major: 1, Minor: 24, Patch: -1, Dev: true}},
	{v: "v14.2.34-dev", exp: Version{Major: 14, Minor: 2, Patch: 34, Dev: true}},
}

func TestVersionString(t *testing.T) {
	for _, tst := range versionStringTests {
		got := tst.exp.String()
		want := tst.v
		if got != want {
			t.Log(tst.exp)
			t.Fatalf("got %q want %q\n", got, want)
		}
	}
}
