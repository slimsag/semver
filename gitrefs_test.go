// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var gitRefsTests = []struct {
	file, mainID, mainName string
	capList                []string
	records                []gitRef
}{
	{
		file:     "testdata/github-azul3d-audio",
		mainID:   "cd95fa968a0fa851547bd65e73e1b385a2dca005",
		mainName: "HEAD",
		capList:  []string{"multi_ack", "thin-pack", "side-band", "side-band-64k", "ofs-delta", "shallow", "no-progress", "include-tag", "multi_ack_detailed", "no-done", "symref=HEAD:refs/heads/master", "agent=git/2:2.1.1+github-607-gfba4028"},
		records: []gitRef{
			{Name: "refs/heads/master", Hash: "cd95fa968a0fa851547bd65e73e1b385a2dca005"},
			{Name: "refs/pull/2/head", Hash: "412511b0e46b31cb4eae7323d3db63acfe60bc08"},
			{Name: "refs/pull/4/head", Hash: "e963a7b43a4c4880ca40110550fe0b247e9691c3"},
			{
				Name: "refs/tags/v1", Hash: "f8d048baeca3571b825c647ce6bdc59f9fbf004f",
				PeeledHash: "630ff3922ec7b8b8a76d0f7e26fa40aa76757a92",
			},
		},
	}, {
		file:     "testdata/github-azul3d-dstarlite",
		mainID:   "214bee1c0789085c0db295e570704122a80067e5",
		mainName: "HEAD",
		capList:  []string{"multi_ack", "thin-pack", "side-band", "side-band-64k", "ofs-delta", "shallow", "no-progress", "include-tag", "multi_ack_detailed", "no-done", "symref=HEAD:refs/heads/master", "agent=git/2:2.1.1+github-607-gfba4028"},
		records: []gitRef{
			{Name: "refs/heads/master", Hash: "214bee1c0789085c0db295e570704122a80067e5"},
			{
				Name: "refs/tags/v1", Hash: "6c8dbd02cac610727c10d365d842218c1aea315e",
				PeeledHash: "214bee1c0789085c0db295e570704122a80067e5",
			},
		},
	}, {
		file:     "testdata/github-azul3d-gfx-gl2",
		mainID:   "3ece9485246bd9c378d408625ec2159d226b8ac8",
		mainName: "HEAD",
		capList:  []string{"multi_ack", "thin-pack", "side-band", "side-band-64k", "ofs-delta", "shallow", "no-progress", "include-tag", "multi_ack_detailed", "no-done", "symref=HEAD:refs/heads/master", "agent=git/2:2.1.1+github-607-gfba4028"},
		records: []gitRef{
			{Name: "refs/heads/master", Hash: "3ece9485246bd9c378d408625ec2159d226b8ac8"},
			{Name: "refs/heads/v3-dev", Hash: "f7808f38206dbb1df3eac8ce38dd5056353b04eb"},
			{Name: "refs/pull/21/head", Hash: "f7808f38206dbb1df3eac8ce38dd5056353b04eb"},
			{Name: "refs/pull/21/merge", Hash: "3b7645798d47da89a7ca02f4421c07ba982bbbca"},
			{
				Name: "refs/tags/v1", Hash: "db26827513f09c4cb5ca2860bfc54deaddd63c40",
				PeeledHash: "7fc8654c35aa31a960cc783f58a0634a5c4c1d99",
			},
			{
				Name: "refs/tags/v2", Hash: "016ef4bb6ad1ce63ba8e5d2325dc3c2c1679f845",
				PeeledHash: "3ece9485246bd9c378d408625ec2159d226b8ac8",
			},
		},
	}, {
		file:     "testdata/github-azul3d-gfx-window",
		mainID:   "daee506ca1b1c5088b1205813397b0e25e2fa9e1",
		mainName: "HEAD",
		capList:  []string{"multi_ack", "thin-pack", "side-band", "side-band-64k", "ofs-delta", "shallow", "no-progress", "include-tag", "multi_ack_detailed", "no-done", "symref=HEAD:refs/heads/master", "agent=git/2:2.1.1+github-607-gfba4028"},
		records: []gitRef{
			{Name: "refs/heads/master", Hash: "daee506ca1b1c5088b1205813397b0e25e2fa9e1"},
			{
				Name: "refs/tags/v1", Hash: "bd84527516dec820711a0123b3a67ad5c2005aff",
				PeeledHash: "043fc03c30fec7f7fd3f456be634b13460d33784",
			},
			{
				Name: "refs/tags/v2", Hash: "c9c5c844c59d169be548768db89113a0e4eae5fa",
				PeeledHash: "daee506ca1b1c5088b1205813397b0e25e2fa9e1",
			},
		},
	}, {
		file:     "testdata/google-code-azul3d",
		mainID:   "227b26555939499162b40a7ab64265e70cd3a790",
		mainName: "HEAD",
		capList:  []string{"multi_ack_detailed", "multi_ack", "side-band-64k", "thin-pack", "ofs-delta", "no-progress", "include-tag", "shallow"},
		records: []gitRef{
			{Name: "refs/heads/master", Hash: "227b26555939499162b40a7ab64265e70cd3a790"},
			{Name: "refs/heads/v0", Hash: "3fcbb5cadc665d0c151d3d042c66ee6c59879b83"},
		},
	},
}

func TestParseRefs(t *testing.T) {
	for _, tst := range gitRefsTests {
		// Read the file data.
		data, err := ioutil.ReadFile(tst.file)
		if err != nil {
			t.Fatal(err)
		}
		dataCpy := make([]byte, len(data))
		copy(dataCpy, data)

		// Parse the refs.
		refs, err := gitParseRefs(data)
		if err != nil {
			t.Log(tst.file)
			t.Fatal(err)
		}

		// Check for a valid service.
		if refs.service != "git-upload-pack" {
			t.Log(tst.file)
			t.Logf("got service=%q\n", refs.service)
			t.Logf(`want service="git-upload-pack"\n`)
			t.Fail()
			continue
		}

		// Check for a valid HEAD.
		if refs.mainID != tst.mainID || refs.mainName != tst.mainName {
			t.Log(tst.file)
			t.Logf("got %q=%q\n", refs.mainName, refs.mainID)
			t.Logf("want %q=%q\n", tst.mainName, tst.mainID)
			t.Fail()
			continue
		}

		// Verify the capabilities list.
		if len(refs.capList) != len(tst.capList) {
			t.Log(tst.file)
			t.Log("got caps", refs.capList)
			t.Log("want caps", tst.capList)
			t.Fail()
			continue
		}
		for i, s := range refs.capList {
			if s != tst.capList[i] {
				t.Log(tst.file)
				t.Log("got caps", refs.capList)
				t.Log("want caps", tst.capList)
				t.Logf("%d. got=%q, want=%q\n", i, s, tst.capList[i])
				t.Fail()
				break
			}
		}

		// Verify the ref records.
		if len(refs.records) != len(tst.records) {
			t.Log(tst.file)
			t.Log("got records", refs.records)
			t.Log("want records", tst.records)
			t.Fail()
			continue
		}
		for i, s := range refs.records {
			if *s != tst.records[i] {
				t.Log(tst.file)
				t.Log("got records", refs.records)
				t.Log("want records", tst.records)
				t.Logf("%d. got=%q, want=%q\n", i, s, tst.records[i])
				t.Fail()
				break
			}
		}

		// Verify that the original data slice was untouched.
		if !bytes.Equal(data, dataCpy) {
			t.Log(tst.file)
			t.Fatal("original data slice is modified")
		}

		// Encode the refs.
		enc := refs.Bytes()
		if !bytes.Equal(enc, data) {
			t.Log(tst.file)
			t.Logf("Want:\n\n%q\n", string(data))
			t.Logf("Got:\n\n%q\n", string(enc))
			t.Fatal("invalid conversion back to bytes")
		}
	}
}
