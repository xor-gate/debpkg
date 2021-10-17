// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xor-gate/debpkg/internal/test"
)

// testWrite writes the deb package to a temporary file and verifies with native dpkg tool when available
func testWrite(t *testing.T, deb *DebPkg) error {
	f := test.TempFile(t)
	err := deb.Write(f)
	if err != nil {
		return err
	}
	err = testReadWithLintian(t, f)
	if err != nil {
		return err
	}
	testReadWithNativeDpkg(t, f)
	return nil
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

	// TODO test control, extract, verify...
	assert.Nil(t, dpkg("info", filename))
	assert.Nil(t, dpkg("contents", filename))
}

func testReadWithLintian(t *testing.T, filename string) error {
	lintianCmd, err := exec.LookPath("lintian")
	if err != nil || lintianCmd == "" {
		return nil
	}

	lintian := func(filename string) error {
		// For now we don't fail on warning or errors (yet)
		//return exec.Command(lintianCmd, "--fail-on", "warning,error", filename).Run()
		return exec.Command(lintianCmd, "--fail-on", "pedantic", filename).Run()
	}

	err = lintian(filename)
	assert.Nil(t, err)
	return err
}

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	defer deb.Close()
	deb.SetName("debpkg-test-add-directory")
	deb.SetArchitecture("all")

	assert.Nil(t, deb.AddDirectory("internal"))
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
	deb.AddFile("debpkg_test.go", "/foo/awesome/test.go")
	deb.AddFileString("this is a real file", "/real/file.txt")

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

// ExampleDebPkgWrite demonstrates generating a simple package
func ExampleDebPkg_Write() {
	tempfile := os.TempDir() + "/foobar.deb"

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

// TestFilenameFromFullVersion verifies if the whole version string is correctly calculated
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
