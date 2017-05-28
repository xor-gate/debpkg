// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
)

const debPkgDigestDefaultHash = crypto.SHA1
const debPkgDigestVersion = 4
const debPkgDigestRole = "builder"

// Digest file for GPG signing
type debPkgDigest struct {
	plaintext string // Plaintext package digest (empty when unsigned)
	clearsign string // GPG clearsigned package digest (empty when unsigned)
	version   int    // Always version 4 (for dpkg-sig 0.13.1+nmu2)
	signer    string // Name <email>
	date      string // Mon Jan 2 15:04:05 2006 (time.ANSIC)
	role      string // builder
	files     string // Multiple "\t<md5sum> <sha1sum> <size> <filename>"
	// E.g:
	//       3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	//       79bb73dbb522dc1a2dd1b9c2ec89fc79 26d29d15aad5c0e051d07571e28da2bc0009707e 366 control.tar.gz
	//       e1a6e48c95a760170029ef7872cec994 e02ed99e5c4fd847bde12b4c2c30dd814b26ec27 136 data.tar.gz
}

// Create unsigned digest file at toplevel of deb package
// NOTE: the deb.digest.version and deb.digest.role are set in this function!
func createDigestFileString(deb *DebPkg) string {
	const digestFileTmpl = `Version: %d
Signer: %s
Date: %s
Role: %s
Files: 
%s`
	deb.digest.version = debPkgDigestVersion
	deb.digest.role = debPkgDigestRole

	// debian-binary
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), md5.New()),
		digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), sha1.New()),
		len(deb.debianBinary),
		"debian-binary")

	// control.tar.gz
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		0, 0,
		len(deb.control.buf.Bytes()),
		"control.tar.gz")

	// data.tar.gz
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		0, 0,
		len(deb.data.buf.Bytes()),
		"data.tar.gz")

	return fmt.Sprintf(digestFileTmpl,
		deb.digest.version,
		deb.digest.signer,
		deb.digest.date,
		deb.digest.role,
		deb.digest.files)
}

func digestCalcDataHash(data *bytes.Buffer, hash hash.Hash) string {
	var result []byte
	if _, err := io.Copy(hash, data); err != nil {
		return ""
	}
	return string(hash.Sum(result))
}
