// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"os"
	"time"

	"github.com/blakesmith/ar"
)

func addArFile(now time.Time, w *ar.Writer, name string, body []byte) error {
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

func (deb *DebPkg) createDebAr(filename string) error {
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

	deb.data.tgz.Close()

	now := time.Now()
	w := ar.NewWriter(fd)

	if err := w.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("cannot write ar header to deb file: %v", err)
	}
	if err := addArFile(now, w, "debian-binary", []byte(deb.debianBinary)); err != nil {
		return fmt.Errorf("cannot pack debian-binary: %v", err)
	}
	if err := addArFile(now, w, "control.tar.gz", deb.control.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add control.tar.gz to deb: %v", err)
	}
	if err := addArFile(now, w, "data.tar.gz", deb.data.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add data.tar.gz to deb: %v", err)
	}
	if deb.digest.clearsign != "" {
		if err := addArFile(now, w, "digests.asc", []byte(deb.digest.clearsign)); err != nil {
			return fmt.Errorf("cannot add digests.asc to deb: %v", err)
		}
	}

	removeDeb = false

	return nil
}
