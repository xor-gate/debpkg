// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xor-gate/debpkg/lib/targzip"
)

type data struct {
	md5sums string
	tgz     *targzip.TarGzip
	dirs    []string
}

func (d *data) addDirectory(dirpath string) error {
	for _, addedDir := range d.dirs {
		if addedDir == dirpath {
			return nil
		}
	}
	if err := d.tgz.AddDirectory(dirpath); err != nil {
		return err
	}
	d.dirs = append(d.dirs, dirpath)
	return nil
}

func (d *data) addEmptyDirectory(dir string) error {
	dirname := strings.Replace(dir, "\\", "/", -1)
	dirs := strings.Split(dirname, "/")
	var current string
	for _, dir := range dirs {
		if len(dir) > 0 {
			current += dir + "/"
			err := d.addDirectory(current)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *data) addFile(filename string, dest ...string) error {
	if err := d.tgz.AddFile(filename, dest...); err != nil {
		return err
	}

	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	md5, err := computeMd5(fd)
	if err != nil {
		fd.Close()
		return err
	}

	d.md5sums += fmt.Sprintf("%x  %s\n", md5, filename)

	fd.Close()
	return nil
}

// computeMd5 from the os filedescriptor
func computeMd5(fd io.Reader) (data []byte, err error) {
	var result []byte
	hash := md5.New()
	if _, err := io.Copy(hash, fd); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}
