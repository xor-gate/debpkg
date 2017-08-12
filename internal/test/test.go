// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package test

import (
	"os"
	"io/ioutil"
	"testing"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

var tempdir string

func init() {
	tempdir, _ = ioutil.TempDir("", "debpkg-test")
}

// TempDir returns the current tempdir
func TempDir() string {
	return tempdir
}

// TempFile calculates the debian package filename based on os.TempDir() and t.Name()
func TempFile(t *testing.T) string {
	return TempDir() + string(os.PathSeparator) + t.Name() + ".deb"
}

// TempOpenPGPIdentity creates a new identity in TempDir()
func TempOpenPGPIdentity() (e *openpgp.Entity, err error) {
	// Create random new GPG identity for signage
	e, _ = openpgp.NewEntity("Debpkg Authors", "", "debpkg-authors@xor-gate.org", nil)

	// Sign all the identities
	for _, id := range e.Identities {
		if err = id.SelfSignature.SignUserId(id.UserId.Id, e.PrimaryKey, e.PrivateKey, nil); err != nil {
			return
		}
	}

	f, _ := os.Create(TempDir() + string(os.PathSeparator) + "openpgp-testkey.asc")
	w, _ := armor.Encode(f, openpgp.PublicKeyType, nil)
	devnull, _ := os.Open(os.DevNull)
	e.SerializePrivate(devnull, nil)
	devnull.Close()
	e.Serialize(w)
	w.Close()
	f.Close()

	return
}
