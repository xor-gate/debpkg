// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"time"

	"github.com/xor-gate/debpkg/lib/targzip"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
)

// DebPkg holds data for a single debian package
type DebPkg struct {
	debianBinary string
	control      debPkgControl
	data         debPkgData
	digest       debPkgDigest
}

// New creates new debian package
func New() *DebPkg {
	deb := &DebPkg{}

	deb.debianBinary = debPkgDebianBinary
	deb.control.info.vcsType = VcsTypeUnset
	deb.control.info.priority = PriorityUnset

	deb.control.buf = &bytes.Buffer{}
	deb.control.tgz = targzip.New(deb.control.buf)

	deb.data.buf = &bytes.Buffer{}
	deb.data.tgz = targzip.New(deb.data.buf)

	return deb
}

func (deb *DebPkg) verify() error {
	if deb.control.info.name == "" {
		return fmt.Errorf("empty package name")
	}

	if deb.control.info.architecture == "" {
		return fmt.Errorf("empty architecture")
	}

	return nil
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	err := deb.verify()
	if err != nil {
		return err
	}

	err = createControlTarGz(deb)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	if filename == "" {
		filename = deb.GetFilename()
	}

	return deb.createDebAr(filename)
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
		debPkgDebianFileExtension)
}

// WriteSigned package with GPG entity
func (deb *DebPkg) WriteSigned(filename string, entity *openpgp.Entity, keyid string) error {
	var buf bytes.Buffer
	var cfg packet.Config
	var signer string
	cfg.DefaultHash = debPkgDigestDefaultHash

	for id := range entity.Identities {
		// TODO real search for keyid, need to investigate maybe a subkey?
		signer = id
	}

	deb.digest.date = fmt.Sprintf(time.Now().Format(time.ANSIC))
	deb.digest.signer = signer

	clearsign, err := clearsign.Encode(&buf, entity.PrivateKey, &cfg)
	if err != nil {
		return fmt.Errorf("error while signing: %s", err)
	}

	err = createControlTarGz(deb)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	deb.digest.plaintext = createDigestFileString(deb)

	if _, err = clearsign.Write([]byte(deb.digest.plaintext)); err != nil {
		return fmt.Errorf("error from Write: %s", err)
	}

	if err = clearsign.Close(); err != nil {
		return fmt.Errorf("error from Close: %s", err)
	}

	deb.digest.clearsign = buf.String()

	return deb.createDebAr(filename)
}

// AddFile adds a file by filename to the package
func (deb *DebPkg) AddFile(filename string, dest ...string) error {
	return deb.data.addFile(filename, dest...)
}

// AddEmptyDirectory adds a empty directory to the package
func (deb *DebPkg) AddEmptyDirectory(dir string) error {
	return deb.data.addEmptyDirectory(dir)
}

// AddDirectory adds a directory to the package
func (deb *DebPkg) AddDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "." || path == ".." {
			return nil
		}

		err = deb.AddFile(path)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetArchitecture gets the current local CPU architecture in debian-form
func GetArchitecture() string {
	arch := build.Default.GOARCH
	if arch == "386" {
		return "i386"
	}
	return arch
}
