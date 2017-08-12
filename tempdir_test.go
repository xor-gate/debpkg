// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

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
	tempDir = ""
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

// TestTempDirNotWritable test if setting a tempdir which is not writeable returns an error
func TestTempDirNotWritable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	assert.Equal(t, os.TempDir(), TempDir())
	assert.NotNil(t, SetTempDir("/this/is/not/writable"))
	assert.Equal(t, os.TempDir(), TempDir())
}

// TestTempDirNewErr
func TestTempDirNewErr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	tempDir = "/this/is/not/writable"

	deb := New()
	defer deb.Close()

	// Deb should contain an I/O error and Write() should also fail
	assert.Equal(t, ErrIO, deb.err)
	assert.Equal(t, ErrIO, deb.Write(""))

	// Restore to os.TempDir()
	assert.Nil(t, SetTempDir(""))
}
