// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xor-gate/debpkg/internal/targzip"
)

// Package holds data of a single pkg.an package
type Package struct {
	Variables Variables
	debianBinary string
	control      Control
	data         data
	digest       digest
	err          error
}

// New creates new debian package, optionally provide an tempdir to write
//  intermediate files, otherwise os.TempDir is used. A provided tempdir must exist
//  in order for it to work.
func New(tempDir ...string) *Package {
	pkg := &Package{
		Variables: DefaultVariables(),
		debianBinary: debianBinaryVersion,
	}

	dir := os.TempDir()
	if len(tempDir) > 0 && len(tempDir[0]) > 0 {
		dir = tempDir[0]
	}

	control, err := targzip.NewTempFile(dir)
	if err != nil {
		pkg.setError(ErrIO)
		return pkg
	}

	data, err := targzip.NewTempFile(dir)
	if err != nil {
		control.Close()
		control.Remove()
		pkg.setError(ErrIO)
		return pkg
	}

	pkg.control.tgz = control
	pkg.data.tgz = data

	return pkg
}

// Close closes the File (and removes the intermediate files), rendering it unusable for I/O. It returns an error, if any.
func (pkg *Package) Close() error {
	if pkg.err == ErrClosed {
		return pkg.err
	}
	if pkg.control.tgz != nil {
		pkg.control.tgz.Remove()
	}
	if pkg.data.tgz != nil {
		pkg.data.tgz.Remove()
	}
	pkg.err = ErrClosed // FIXME make pkg.SetError work...
	return nil
}

// writeControlData writes the control.tar.gz
func (pkg *Package) writeControlData() error {
	err := pkg.control.Verify()
	if err != nil {
		return err
	}

	err = pkg.control.finalizeControlFile(&pkg.data)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	if err := pkg.control.tgz.Close(); err != nil {
		return fmt.Errorf("cannot close tgz writer: %v", err)
	}

	if err := pkg.data.tgz.Close(); err != nil {
		return fmt.Errorf("cannot close tgz writer: %v", err)
	}
	return nil
}

// Write the debian package to the filename
func (pkg *Package) Write(filename string) error {
	if pkg.err != nil {
		return pkg.err
	}
	if err := pkg.writeControlData(); err != nil {
		pkg.setError(err)
		return err
	}
	if filename == "" {
		filename = pkg.GetFilename()
	}
	err := pkg.createDebAr(filename)
	pkg.setError(err)
	pkg.Close()
	return err
}

// GetFilename calculates the filename based on name, version and architecture
// SetName("foo")
// SetVersion("1.33.7")
// SetArchitecture("amd64")
// Generates filename "foo-1.33.7_amd64.deb"
func (pkg *Package) GetFilename() string {
	return fmt.Sprintf("%s-%s_%s.%s",
		pkg.control.info.name,
		pkg.control.info.version.Full,
		pkg.control.info.architecture,
		debianFileExtension)
}

// MarkConfigFile marks configuration files in the pkg.an package
func (pkg *Package) MarkConfigFile(dest string) error {
	return pkg.control.markConfigFile(dest)
}

// AddFile adds a file by filename to the package
func (pkg *Package) AddFile(filename string, dest ...string) error {
	if pkg.err != nil {
		return pkg.err
	}
	return pkg.setError(pkg.data.addFile(filename, dest...))
}

// AddFileString adds a file to the package with the provided content
func (pkg *Package) AddFileString(contents, dest string) error {
	if pkg.err != nil {
		return pkg.err
	}
	return pkg.setError(pkg.data.addFileString(contents, dest))
}

// AddEmptyDirectory adds a empty directory to the package
func (pkg *Package) AddEmptyDirectory(dir string) error {
	if pkg.err != nil {
		return pkg.err
	}
	return pkg.setError(pkg.data.addDirectory(dir))
}

// AddDirectory adds a directory recursive to the package
func (pkg *Package) AddDirectory(dir string) error {
	if pkg.err != nil {
		return pkg.err
	}

	pkg.data.addDirectory(dir)

	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." || path == ".." || dir == path {
			return nil
		}
		if f.IsDir() {
			if err := pkg.data.addDirectory(path); err != nil {
				return pkg.setError(err)
			}
			return pkg.AddDirectory(path)
		}

		return pkg.setError(pkg.AddFile(path))
	})
}
