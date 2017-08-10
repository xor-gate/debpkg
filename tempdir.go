// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"os"
)

var tempDir = os.TempDir() // default temporary directory is os.TempDir

// SetTempDir sets the directory for temporary files. When the directory doesn't
//  exist it is automaticly created (but not removed).
func SetTempDir(dir string) error {
	if dir == "" {
		dir = os.TempDir()
	}

	stat, err := os.Stat(dir)
	if err == nil && stat.IsDir() {
		tempDir = dir
		return nil
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	tempDir = dir
	return nil
}

// RemoveTempDir removes the temporary directory recursive. This is safe against
//  when TempDir() is set to os.TempDir() then it does nothing
func RemoveTempDir() error {
	if TempDir() == os.TempDir() {
		return nil
	}
	return os.RemoveAll(TempDir())
}

// TempDir returns the directory to use for temporary files.
func TempDir() string {
	return tempDir
}
