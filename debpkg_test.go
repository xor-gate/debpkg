// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	defer deb.Close()
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

func TestWrite(t *testing.T) {
	deb := New()
	defer deb.Close()
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
		t.Errorf("Error in writing package: %v", err)
	}
}

func TestWriteError(t *testing.T) {
	deb := New()
	defer deb.Close()
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
	defer deb.Close()

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
	return exec.Command(cmd, "--"+action, filename).Run()
}

func TestReadWithNativeDpkg(t *testing.T) {
	dpkgCmd, err := exec.LookPath("dpkg")
	if err != nil || dpkgCmd == "" {
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

		dpkg(dpkgCmd, "contents", deb)
		if err != nil {
			t.Errorf("dpkg --contents failed on " + deb)
		}
	}
}

func TestFilenameFromFullVersion(t *testing.T) {
	deb := New()
	defer deb.Close()

	deb.SetName("foo")
	deb.SetVersion("1.33.7")
	deb.SetArchitecture("amd64")

	fn := deb.GetFilename()
	if fn != "foo-1.33.7_amd64.deb" {
		t.Errorf("unexpected filename: " + fn)
	}
}
