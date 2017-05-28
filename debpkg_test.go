// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/openpgp"
)

var e *openpgp.Entity

func init() {
	// Create random new GPG identity for signage
	e, _ = openpgp.NewEntity("Foo Bar", "", "foo@bar.com", nil)
}

// Test creation of empty digest
func TestDigestCreateEmpty(t *testing.T) {
	// FIXME it seems whe digesting the data buf the whole tarball will go corrupt...
	/*
	   	digestExpect := `Version: 4
	   Signer:
	   Date:
	   Role: builder
	   Files:
	   	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	   	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 control.tar.gz
	   	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 data.tar.gz
	   `
	*/
	digestExpect := `Version: 4
Signer: 
Date: 
Role: builder
Files: 
	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	0 0 0 control.tar.gz
	0 0 0 data.tar.gz
`

	deb := New()
	digest := createDigestFileString(deb)

	if digest != digestExpect {
		t.Error("Unexpected digest file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(digestExpect), digestExpect, len(digest), digest)
	}
}

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	deb.SetName("debpkg-test-add-directory")
	deb.SetArchitecture("all")
	err := deb.AddDirectory("vendor")
	if err != nil {
		t.Errorf("Error adding directory '.': %v", err)
		return
	}

	err = deb.Write("debpkg-test-add-directory.deb")
	if err != nil {
		t.Errorf("Error writing debfile: %v", err)
		return
	}
}

func TestWriteSignedEmpty(t *testing.T) {
	deb := New()

	// WriteSigned package
	err := deb.WriteSigned("debpkg-test-signed-empty.deb", e, "00000000")
	if err != nil {
		t.Errorf("Error in writing signed package: %v", err)
	}
}

func TestWrite(t *testing.T) {
	deb := New()

	deb.SetName("debpkg-test")
	deb.SetArchitecture("all")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://foobar.com")
	deb.SetShortDescription("some awesome foobar pkg")
	deb.SetDescription("very very very very long description")

	// Set version control system info for control file
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/secdl")
	deb.SetVcsBrowser("https://github.com/xor-gate/secdl")

	deb.AddFile("debpkg.go")

	err := deb.Write("debpkg-test.deb")
	if err != nil {
		t.Errorf("Error in writing unsigned package: %v", err)
	}
}

func TestWriteSigned(t *testing.T) {
	deb := New()

	deb.SetName("debpkg-test-signed")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
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

func TestWriteError(t *testing.T) {
	deb := New()
	err := deb.Write("")
	if err == nil {
		t.Errorf("deb.Write shouldnt return nil")
	}
	deb.control.info.name = "pkg"
	if err := deb.Write(""); err == nil {
		t.Errorf("deb.Write shouldnt return nil")
	}
}

func ExampleWrite() {
	deb := New()

	deb.SetName("foobar")
	deb.SetVersion("1.2.3")
	deb.SetArchitecture("amd64")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("http://foobar.com")

	deb.SetShortDescription("Minimal foo bar package")
	deb.SetDescription("Foo bar package doesn't do anything")

	deb.AddFile("debpkg.go")
	fmt.Println(deb.Write("foobar.deb"))

	// Output: <nil>
}

func dpkg(cmd, action, filename string) error {
	args := []string{"--" + action, filename}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return err
	}
	return nil
}

func TestReadWithNativeDpkg(t *testing.T) {
	dpkgCmd, err := exec.LookPath("dpkg")
	if err != nil || dpkgCmd == "" {
		fmt.Println("Skip test, unable to find dpkg in PATH")
		return
	}

	debs, err := filepath.Glob("*.deb")
	if err != nil {
		t.Errorf("Unexpected error on glob: %v", err)
	}
	for _, deb := range debs {
		err = dpkg(dpkgCmd, "info", deb)
		if err != nil {
			t.Errorf("dpkg --info failed on " + deb)
		}
		fmt.Println("dpkg --info passed on " + deb)

		dpkg(dpkgCmd, "contents", deb)
		if err != nil {
			t.Errorf("dpkg --contents failed on " + deb)
		}
		fmt.Println("dpkg --contents passed on " + deb)
	}
}

func TestFilenameFromFullVersion(t *testing.T) {
	deb := New()

	deb.SetName("foo")
	deb.SetVersion("1.33.7")
	deb.SetArchitecture("amd64")

	fn := deb.GetFilename()
	if fn != "foo-1.33.7_amd64.deb" {
		t.Errorf("unexpected filename: " + fn)
	}
}
