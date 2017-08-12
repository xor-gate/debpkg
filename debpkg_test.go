// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"go/build"
	"os/exec"
	"testing"
	"os"
	"github.com/xor-gate/debpkg/internal/test"
	"github.com/stretchr/testify/assert"
)

// testWrite writes the deb package to a temporary file
func testWrite(t *testing.T, deb *DebPkg) error {
	f := test.TempFile(t)
	err := deb.Write(f)
	if err == nil {
		testReadWithNativeDpkg(t, f)
	}
	return err
}

// testReadWithNativeDpkg tests a single debian package with the dpkg tool when present
func testReadWithNativeDpkg(t *testing.T, filename string) {
	dpkgCmd, err := exec.LookPath("dpkg")
	if err != nil || dpkgCmd == "" {
		return
	}

	dpkg := func(action, filename string) error {
		return exec.Command(dpkgCmd, "--"+action, filename).Run()
	}

	// TODO test dry-run install...
	assert.Nil(t, dpkg("info", filename))
	assert.Nil(t, dpkg("contents", filename))
}

// TestTempDir verifies the correct working of TempDir and SetTempDir
func TestTempDir(t *testing.T) {
	dirExists := func(path string) bool {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return true
		}
		return false
	}

	// Default the TempDir points to os.TempDir()
	assert.Equal(t, os.TempDir(), TempDir())

	// Unset debpkgTempDir and verify it is set to os.TempDir() when SetTempDir received a empty string
	debpkgTempDir = ""
	assert.Nil(t, SetTempDir(""))
	assert.Equal(t, os.TempDir(), TempDir())

	// Check if custom test tempdir is created
	tempdir := os.TempDir() + "/debpkg-test-tempdir"

	assert.Nil(t, SetTempDir(tempdir))
	assert.True(t, dirExists(tempdir))
	assert.Nil(t, RemoveTempDir())
	assert.False(t, dirExists(tempdir))
	assert.Nil(t, SetTempDir(""))

	// Check if TempDir() == os.TempDir() is not removed and RemoveTempDir() returns nil on os.TempDir()
	assert.True(t, dirExists(TempDir()))
	assert.Nil(t, RemoveTempDir())
	assert.True(t, dirExists(TempDir()))

	// Restore to os.TempDir()
	assert.Nil(t, SetTempDir(""))
}

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	defer deb.Close()
	deb.SetName("debpkg-test-add-directory")
	deb.SetArchitecture("all")

	assert.Nil(t, deb.AddDirectory("vendor"))
	assert.Nil(t, testWrite(t, deb))
}

// TestWrite verifies Write works as expected with adding just one datafile
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

	assert.Nil(t, testWrite(t, deb))

	// Try to Write again on implicit closed package
	assert.Equal(t, ErrClosed, testWrite(t, deb))
}

// TestWriteError tests if the Write fails with the correct errors
func TestWriteError(t *testing.T) {
	deb := New()
	defer deb.Close()
	assert.NotNil(t, deb.Write(""), "deb.Write should return nil")

	deb.control.info.name = "pkg"
	assert.Equal(t, fmt.Errorf("empty package name"), deb.Write(""))
}

// ExampleWrite demonstrates generating a simple package
func ExampleWrite() {
	tempfile := TempDir() + "/foobar.deb"

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
	fmt.Println(deb.Write(tempfile))

	// Do something with tempfile other than removing it...
}

// TestFilenameFromFullVersion verifies if the whole version string is correctly concatinated
func TestFilenameFromFullVersion(t *testing.T) {
	deb := New()
	defer deb.Close()

	deb.SetName("foo")
	deb.SetVersion("1.33.7")
	deb.SetArchitecture("amd64")

	assert.Equal(t, "foo-1.33.7_amd64.deb", deb.GetFilename())
}

// TestGetArchitecture checks the current build.Default.GOARCH compatible debian architecture
func TestGetArchitecture(t *testing.T) {
	// On debian 386 GOARCH is presented as i386
	goarch := build.Default.GOARCH
	build.Default.GOARCH = "386"
	assert.Equal(t, "i386", GetArchitecture())
	build.Default.GOARCH = goarch

	// Check current build GOARCH
	if build.Default.GOARCH != "386" {
		assert.Equal(t, build.Default.GOARCH, GetArchitecture())
	}
}
