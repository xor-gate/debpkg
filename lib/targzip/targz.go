// Copyright 2017 Jerry Jacobs. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package targzip

import (
	"os"
	"io"
	"fmt"
	"time"
	"strings"
	"archive/tar"
	"path/filepath"
	"compress/gzip"
)

type TarGzip struct {
	tw *tar.Writer
	gw *gzip.Writer
}

// New creates a new targzip writer
func New(w io.Writer) *TarGzip {
	t := &TarGzip{}

	t.gw = gzip.NewWriter(w)
	t.tw = tar.NewWriter(t.gw)

	return t
}

// AddFile write a file from filename into dest
func (t *TarGzip) AddFile(filename string, dest ...string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		return err
	}
	if stat.Mode().IsDir() {
		return nil
	}

	dirname := filepath.Dir(filename)
	if dirname != "." {
		dirname = strings.Replace(dirname, "\\", "/", -1)
		dirs := strings.Split(dirname, "/")
		var current string
		for _, dir := range dirs {
			if len(dir) > 0 {
				current += dir + "/"
				t.AddDirectory(current)
			}
		}
	}

	// now lets create the header as needed for this file within the tarball
	hdr, err := tar.FileInfoHeader(stat, filename)
	if err != nil {
		return fmt.Errorf("dir tar finfo: %v", err)
	}
	if len(dest) > 0 {
		hdr.Name = dest[0]
	} else {
		hdr.Name = filename
	}

	// write the header to the tarball archive
	if err := t.WriteHeader(hdr); err != nil {
		return err
	}

	// copy the file data to the tarball
	if _, err := io.Copy(t, fd); err != nil {
		return err
	}

	return nil
}

// AddFileFromBuffer adds a file from a buffer
func (t *TarGzip) AddFileFromBuffer(filename string, b []byte) error {
	hdr := tar.Header{
		Name:     filename,
		Size:     int64(len(b)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if err := t.WriteHeader(&hdr); err != nil {
		return fmt.Errorf("cannot write header of file: %v", err)
	}

	if _, err := t.Write(b); err != nil {
		return fmt.Errorf("cannot write file: %v", err)
	}

	return nil
}

// AddDirectory adds a directory entry
func (t *TarGzip) AddDirectory(dirpath string) error {
	hdr := &tar.Header{
		Name:     dirpath,
		Mode:     int64(0755 | 040000),
		Typeflag: tar.TypeDir,
		ModTime:  time.Now(),
		Size:     0,
	}

	if err := t.WriteHeader(hdr); err != nil {
		return fmt.Errorf("tar-header for dir: %v", err)
	}

	return nil
}

// WriteHeader writes a raw tar header
func (t *TarGzip) WriteHeader(hdr *tar.Header) error {
	return t.tw.WriteHeader(hdr)
}

// Write writes raw tar data
func (t *TarGzip) Write(p []byte) (n int, err error) {
	return t.tw.Write(p)
}

// Close closes the targzip writer
func (t *TarGzip) Close() error {
	if err := t.tw.Close(); err != nil {
		return err
	}
	if err := t.gw.Close(); err != nil {
		return err
	}
	return nil
}
