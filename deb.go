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

// DebPkg holds data for a single debian package
type DebPkg struct {
	debianBinary string
	control      control
	data         data
	digest       digest
	err          error
}

// New creates new debian package
func New() *DebPkg {
	deb := &DebPkg{
		debianBinary: debianBinaryVersion,
	}

	control, err := targzip.NewTempFile(tempDir)
	if err != nil {
		deb.setError(ErrIO)
		return deb
	}

	data, err := targzip.NewTempFile(tempDir)
	if err != nil {
		control.Close()
		control.Remove()
		deb.setError(ErrIO)
		return deb
	}

	deb.control.tgz = control
	deb.data.tgz = data

	return deb
}

// Close closes the File (and removes the intermediate files), rendering it unusable for I/O. It returns an error, if any.
func (deb *DebPkg) Close() error {
	if deb.err == ErrClosed {
		return deb.err
	}
	if deb.control.tgz != nil {
		deb.control.tgz.Remove()
	}
	if deb.data.tgz != nil {
		deb.data.tgz.Remove()
	}
	deb.err = ErrClosed // FIXME make deb.SetError work...
	return nil
}

// writeControlData writes the control.tar.gz
func (deb *DebPkg) writeControlData() error {
	err := deb.control.verify()
	if err != nil {
		return err
	}

	err = createControlTarGz(deb)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	if err := deb.control.tgz.Close(); err != nil {
		return fmt.Errorf("cannot close tgz writer: %v", err)
	}

	if err := deb.data.tgz.Close(); err != nil {
		return fmt.Errorf("cannot close tgz writer: %v", err)
	}
	return nil
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	if deb.err != nil {
		return deb.err
	}
	if err := deb.writeControlData(); err != nil {
		deb.setError(err)
		return err
	}
	if filename == "" {
		filename = deb.GetFilename()
	}
	err := deb.createDebAr(filename)
	deb.setError(err)
	deb.Close()
	return err
}

// GetFilename calculates the filename based on name, version and architecture
// SetName("foo")
// SetVersion("1.33.7")
// SetArchitecture("amd64")
// Generates filename "foo-1.33.7_amd64.deb"
func (deb *DebPkg) GetFilename() string {
	return fmt.Sprintf("%s-%s_%s.%s",
		deb.control.info.name,
		deb.control.info.version.full,
		deb.control.info.architecture,
		debianFileExtension)
}

// AddFile adds a file by filename to the package
func (deb *DebPkg) AddFile(filename string, dest ...string) error {
	if deb.err != nil {
		return deb.err
	}
	return deb.setError(deb.data.addFile(filename, dest...))
}

// AddFileString adds a file to the package with the provided content
func (deb *DebPkg) AddFileString(contents, dest string) error {
	if deb.err != nil {
		return deb.err
	}
	return deb.setError(deb.data.addFileString(contents, dest))
}

// AddEmptyDirectory adds a empty directory to the package
func (deb *DebPkg) AddEmptyDirectory(dir string) error {
	if deb.err != nil {
		return deb.err
	}
	return deb.setError(deb.data.addDirectory(dir))
}

// AddDirectory adds a directory recursive to the package
func (deb *DebPkg) AddDirectory(dir string) error {
	if deb.err != nil {
		return deb.err
	}

	deb.data.addDirectory(dir)

	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." || path == ".." || dir == path {
			return nil
		}
		if f.IsDir() {
			if err := deb.data.addDirectory(path); err != nil {
				return deb.setError(err)
			}
			return deb.AddDirectory(path)
		}

		return deb.setError(deb.AddFile(path))
	})
}
