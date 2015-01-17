// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"testing"
)

var gitBasicPktLinesTests = []struct {
	encoded, decoded string
}{
	{encoded: "0006a\n", decoded: "a\n"},
	{encoded: "0005a", decoded: "a"},
	{encoded: "000bfoobar\n", decoded: "foobar\n"},
	{encoded: "0004", decoded: ""},
}

func TestGitBasicPktLines(t *testing.T) {
	for _, tst := range gitBasicPktLinesTests {
		// Parse the pkt-line.
		pl, _, _, err := gitNextPktLine([]byte(tst.encoded))
		if err != nil {
			t.Logf("encoded %q\n", tst.encoded)
			t.Fatal(err)
		}

		// Test if the decoded version is correct.
		if tst.decoded != string(pl) {
			t.Logf("encoded %q\n", tst.encoded)
			t.Logf("got     %q\n", pl)
			t.Logf("want    %q\n", tst.decoded)
			t.Fatal(err)
		}

		// Test if re-encoding produces an identical string.
		enc := string(pl.Bytes())
		if tst.encoded != enc {
			t.Logf("encoded %q\n", tst.encoded)
			t.Logf("got     %q\n", enc)
			t.Logf("want    %q\n", tst.encoded)
			t.Fatal(err)
		}
	}
}

func TestGitPktLineBreak(t *testing.T) {
	// Special case: "0000" should return err=gitPktLineBreak
	_, lineBreak, _, err := gitNextPktLine([]byte("0000"))
	if err != nil {
		t.Fatal(err)
	}
	if !lineBreak {
		t.Log(`encoded "0000"`)
		t.Log("got lineBreak=false")
		t.Log("want lineBreak=true")
		t.Fail()
	}
}

func TestGitPktLineStream(t *testing.T) {
	// Test a pkt-line data stream with mixed newlines in it.
	stream := []byte("0006a\n0005a0004000bfoobar\n0000")
	results := []string{
		"a\n", "a", "", "foobar\n", "",
	}
	i := 0
	for {
		if i > len(results)*2 {
			t.Fatal("loop exceeded expected results by 2x")
		}

		// Decode the next pkt-line from the stream.
		pl, lineBreak, n, err := gitNextPktLine(stream)
		if err == errGitPktLineNeedMore {
			// End of stream. Validate that we got the correct number of
			// results.
			if i != len(results) {
				t.Fatalf("got n=%d results, wanted n=%d\n", i, len(results))
			}
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		stream = stream[n:]

		// Validate the decoded pkt-line.
		if string(pl) != results[i] {
			t.Fatalf("%d. got=%q want=%q\n", i, string(pl), results[i])
		}

		// Last result needs to be lineBreak=true, others do not.
		if lineBreak && i < len(results)-1 {
			t.Fatalf("%d. got lineBreak=true, wanted lineBreak=false\n", i)
		} else if !lineBreak && i == len(results)-1 {
			t.Fatalf("%d. got lineBreak=false, wanted lineBreak=true\n", i)
		}
		i++
	}
}
