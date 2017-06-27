// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"os"
	"time"
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
)

const digestDefaultHash = crypto.SHA1
const digestVersion = 4
const digestRole = "builder"

// Digest file for GPG signing
type digest struct {
	plaintext string // Plaintext package digest (empty when unsigned)
	clearsign string // GPG clearsigned package digest (empty when unsigned)
	signer    string // Name <email>
	date      string // Mon Jan 2 15:04:05 2006 (time.ANSIC)
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
	// debian-binary
	md5sum, _ := digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), md5.New())
	sha1sum, _ := digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), sha1.New())
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		md5sum,
		sha1sum,
		len(deb.debianBinary),
		"debian-binary")

	// control.tar.gz
	md5sum, _ = digestCalcDataHashFromFile(deb.control.tgz.Name(), md5.New())
	sha1sum, _ = digestCalcDataHashFromFile(deb.control.tgz.Name(), sha1.New())
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		md5sum,
		sha1sum,
		deb.control.tgz.Size(),
		"control.tar.gz")

	// data.tar.gz
	md5sum, _ = digestCalcDataHashFromFile(deb.data.tgz.Name(), md5.New())
	sha1sum, _ = digestCalcDataHashFromFile(deb.data.tgz.Name(), sha1.New())
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		md5sum,
		sha1sum,
		deb.data.tgz.Size(),
		"data.tar.gz")

	return fmt.Sprintf(digestFileTmpl,
		digestVersion,
		deb.digest.signer,
		deb.digest.date,
		digestRole,
		deb.digest.files)
}

func digestCalcDataHashFromFile(filename string, hash hash.Hash) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return digestCalcDataHash(f, hash)
}

func digestCalcDataHash(in io.Reader, hash hash.Hash) (string, error) {
	var result []byte
	if _, err := io.Copy(hash, in); err != nil {
		return "", err
	}
	return string(hash.Sum(result)),nil
}

// WriteSigned package with GPG entity
func (deb *DebPkg) WriteSigned(filename string, entity *openpgp.Entity) error {
	var buf bytes.Buffer
	var cfg packet.Config
	var signer string
	cfg.DefaultHash = digestDefaultHash

	for id := range entity.Identities {
		// TODO real search for keyid, need to investigate maybe a subkey?
		signer = id
	}

	deb.digest.date = time.Now().Format(time.ANSIC)
	deb.digest.signer = signer

	clearsign, err := clearsign.Encode(&buf, entity.PrivateKey, &cfg)
	if err != nil {
		return fmt.Errorf("error while signing: %s", err)
	}

	if err := deb.writeControlData(); err != nil {
		return err
	}

	deb.digest.plaintext = createDigestFileString(deb)

	if _, err = clearsign.Write([]byte(deb.digest.plaintext)); err != nil {
		return fmt.Errorf("error from Write: %s", err)
	}

	if err = clearsign.Close(); err != nil {
		return fmt.Errorf("error from Close: %s", err)
	}

	deb.digest.clearsign = buf.String()

	if filename == "" {
		filename = deb.GetFilename()
	}
	return deb.createDebAr(filename)
}


