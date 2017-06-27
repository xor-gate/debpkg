// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	defer deb.Close()
	deb.SetName("debpkg-test-add-directory")
	deb.SetArchitecture("all")

	assert.Nil(t, deb.AddDirectory("vendor"))
	assert.Nil(t, deb.Write(os.TempDir() + "/debpkg-test-add-directory.deb"))
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

	assert.Nil(t, deb.Write(os.TempDir() + "/debpkg-test.deb"))
}

func TestWriteError(t *testing.T) {
	deb := New()
	defer deb.Close()
	assert.NotNil(t, deb.Write(""), "deb.Write should return nil")

	deb.control.info.name = "pkg"
	assert.Equal(t, fmt.Errorf("empty architecture"), deb.Write(""))
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
	fmt.Println(deb.Write(os.TempDir() + "/foobar.deb"))

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
	assert.Nil(t, err)
	for _, deb := range debs {
		assert.Nil(t, dpkg(dpkgCmd, "info", deb))
		assert.Nil(t, dpkg(dpkgCmd, "contents", deb))
	}
}

func TestFilenameFromFullVersion(t *testing.T) {
	deb := New()
	defer deb.Close()

	deb.SetName("foo")
	deb.SetVersion("1.33.7")
	deb.SetArchitecture("amd64")

	assert.Equal(t, "foo-1.33.7_amd64.deb", deb.GetFilename())
}
