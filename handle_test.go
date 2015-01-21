// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import "testing"

var refTestData = map[string]*gitRef{
	"v0": &gitRef{
		Name: "refs/heads/master",
		Hash: "001",
	},
	"v1": &gitRef{
		Name:       "refs/tags/v1",
		Hash:       "002",
		PeeledHash: "003",
	},
	"v1.0.1": &gitRef{
		Name: "refs/tags/v1.0.1",
		Hash: "004",
	},
	"v1.2": &gitRef{
		Name:       "refs/tags/v1.2",
		Hash:       "005",
		PeeledHash: "006",
	},
	"v2-unstable": &gitRef{
		Name: "refs/heads/v2-unstable",
		Hash: "007",
	},
}

func testChooseRef(t *testing.T, expect, target string, all []*gitRef) {
	v := ParseVersion(target)
	want := refTestData[expect]
	h := &Handler{}
	chosenHash, ok := h.chooseRef(all, v)
	wantOk := len(expect) > 0
	if ok != wantOk {
		t.Fatalf("chooseRef returned ok=%t, want ok=%t\n", ok, wantOk)
		return
	}
	if !wantOk {
		return
	}
	if chosenHash != want.BestHash() {
		t.Logf("got %q\n", chosenHash)
		t.Fatalf("expected %q\n", want.BestHash())
	}
}

func TestChooseRefAscending(t *testing.T) {
	target := "v1"
	expect := "v1.2"
	testChooseRef(t, expect, target, []*gitRef{
		refTestData["v0"],
		refTestData["v1.0.1"],
		refTestData["v1"],
		refTestData["v1.2"],
		refTestData["v2-unstable"],
	})
}

func TestChooseRefDescending(t *testing.T) {
	target := "v1"
	expect := "v1.2"
	testChooseRef(t, expect, target, []*gitRef{
		refTestData["v2-unstable"],
		refTestData["v1.2"],
		refTestData["v1"],
		refTestData["v1.0.1"],
		refTestData["v0"],
	})
}

func TestChooseRefRandom(t *testing.T) {
	target := "v1"
	expect := "v1.2"
	testChooseRef(t, expect, target, []*gitRef{
		refTestData["v1.0.1"],
		refTestData["v1"],
		refTestData["v2-unstable"],
		refTestData["v1.2"],
		refTestData["v0"],
	})
}

func TestChooseRefUnstable(t *testing.T) {
	target := "v2-unstable"
	expect := "v2-unstable"
	testChooseRef(t, expect, target, []*gitRef{
		refTestData["v1.0.1"],
		refTestData["v1"],
		refTestData["v2-unstable"],
		refTestData["v1.2"],
		refTestData["v0"],
	})
}

func TestChooseRefInvalid(t *testing.T) {
	target := "v2"
	expect := ""
	testChooseRef(t, expect, target, []*gitRef{
		refTestData["v1.0.1"],
		refTestData["v1"],
		refTestData["v2-unstable"],
		refTestData["v1.2"],
		refTestData["v0"],
	})
}
