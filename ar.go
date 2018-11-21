// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/xor-gate/ar"
)

func addArFileFromBuffer(now time.Time, w *ar.Writer, name string, body []byte) error {
	hdr := ar.Header{
		Name:    name,
		Size:    int64(len(body)),
		Mode:    0644,
		ModTime: now,
	}

	if err := w.WriteHeader(&hdr); err != nil {
		return fmt.Errorf("cannot write file header: %v", err)
	}

	_, err := w.Write(body)

	return err
}

func addArFile(now time.Time, w *ar.Writer, dstname, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	hdr := ar.Header{
		Name:    dstname,
		Size:    stat.Size(),
		Mode:    0644,
		ModTime: now,
	}

	if err := w.WriteHeader(&hdr); err != nil {
		f.Close()
		return fmt.Errorf("cannot write file header: %v", err)
	}

	_, err = io.Copy(w, f)
	f.Close()

	return err
}

func (deb *Package) createDebAr(filename string) error {
	removeDeb := true
	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create: %s", filename)
	}

	defer func() {
		fd.Close()
		if removeDeb {
			os.Remove(filename)
		}
	}()

	now := time.Now()
	w := ar.NewWriter(fd)

	if err := w.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("cannot write ar header to deb file: %v", err)
	}
	if err := addArFileFromBuffer(now, w, "debian-binary", []byte(debianBinaryVersion)); err != nil {
		return fmt.Errorf("cannot pack debian-binary: %v", err)
	}
	if err := addArFile(now, w, "control.tar.gz", deb.control.tgz.Name()); err != nil {
		return fmt.Errorf("cannot add control.tar.gz to deb: %v", err)
	}
	if err := addArFile(now, w, "data.tar.gz", deb.data.tgz.Name()); err != nil {
		return fmt.Errorf("cannot add data.tar.gz to deb: %v", err)
	}
	if deb.digest.clearsign != "" {
		if err := addArFileFromBuffer(now, w, "digests.asc", []byte(deb.digest.clearsign)); err != nil {
			return fmt.Errorf("cannot add digests.asc to deb: %v", err)
		}
	}

	removeDeb = false

	return nil
}
