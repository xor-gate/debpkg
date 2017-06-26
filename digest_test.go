// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
//	"fmt"
	"testing"

	"golang.org/x/crypto/openpgp"
)

var e *openpgp.Entity

func init() {
	// Create random new GPG identity for signage
	e, _ = openpgp.NewEntity("Foo Bar", "", "foo@bar.com", nil)
}

/*
func TestDigestCreateEmpty(t *testing.T) {
	digestExpect := `Version: 4
Signer:
Date:
Role: builder
Files:
	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 control.tar.gz
   	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 data.tar.gz
`

	deb := New()
	defer deb.Close()
	digest := createDigestFileString(deb)

	if digest != digestExpect {
		t.Error("Unexpected digest file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(digestExpect), digestExpect, len(digest), digest)
	}
}
*/

func TestWriteSigned(t *testing.T) {
	deb := New()
	defer deb.Close()

	deb.SetName("debpkg-test-signed")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
	deb.SetArchitecture("any")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://foobar.com")
	deb.SetShortDescription("some awesome foobar pkg")
	deb.SetDescription("very very very very long description")

	// Set version control system info for control file
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/secdl")
	deb.SetVcsBrowser("https://github.com/xor-gate/secdl")
	deb.SetPriority(PriorityRequired)
	deb.SetConflicts("bash")
	deb.SetProvides("boembats")

	deb.AddFile("debpkg.go")

	// WriteSigned the package
	err := deb.WriteSigned("debpkg-test-signed.deb", e, "00000000")
	if err != nil {
		t.Errorf("Error in writing unsigned package: %v", err)
	}
}
