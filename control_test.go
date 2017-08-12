// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/xor-gate/debpkg/internal/test"
	"testing"
)

// Test correct output of a empty control file when no DepPkg Set* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileEmpty(t *testing.T) {
	controlExpect := `Package: 
Version: 0.0.0
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Description: 
`
	// Empty
	deb := New()
	defer deb.Close()

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")

	assert.Equal(t, controlExpect, deb.control.String(0))
}

// Test correct output of a control file when SetVcs* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileVcsAndVcsBrowserFields(t *testing.T) {
	controlExpect := `Package: 
Version: 0.0.0
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Vcs-Git: https://github.com/xor-gate/debpkg.git
Vcs-Browser: https://github.com/xor-gate/debpkg
Description: 
`
	// Empty
	deb := New()
	defer deb.Close()

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/debpkg.git")
	deb.SetVcsBrowser("https://github.com/xor-gate/debpkg")

	assert.Equal(t, controlExpect, deb.control.String(0))
}

// Test correct output of the control file when SetVersion* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileSetVersionMajorMinorPatch(t *testing.T) {
	// Empty
	deb := New()
	defer deb.Close()

	deb.SetName("foobar")
	deb.SetArchitecture("amd64")

	// Set major.minor.patch, leave full version string untouched
	deb.SetVersionMajor(1)
	deb.SetVersionMinor(2)
	deb.SetVersionPatch(3)

	controlExpect := `Package: foobar
Version: 1.2.3
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Description: 
`
	assert.Equal(t, controlExpect, deb.control.String(0))

	// Set full version string, this will overwrite the set SetVersion{Major,Minor,Patch} string
	deb.SetVersion("7.8.9")

	controlExpectFullVersion := `Package: foobar
Version: 7.8.9
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Description: 
`

	assert.Equal(t, controlExpectFullVersion, deb.control.String(0))
}

// Test correct output of control file when the mandatory DepPkg Set* functions are called
// This checks if the long description is formatted according to the debian policy
func TestControlFileLongDescriptionFormatting(t *testing.T) {
	controlExpect := `Package: debpkg
Version: 0.0.0
Architecture: amd64
Maintainer: Jerry Jacobs <foo@bar.com>
Installed-Size: 0
Homepage: https://github.com/xor-gate/debpkg
Description: Golang package for creating (gpg signed) debian packages
 **Features**
 
 * Create simple debian packages from files and folders
 * Add custom control files (preinst, postinst, prerm, postrm etcetera)
 * dpkg like tool with a subset of commands (--contents, --control, --extract, --info)
 * Create package from debpkg.yml specfile (like packager.io without cruft)
 * GPG sign package
 * GPG verify package`

	// User supplied very long description without leading spaces and no ending newline
	controlDescr := `**Features**

* Create simple debian packages from files and folders
* Add custom control files (preinst, postinst, prerm, postrm etcetera)
* dpkg like tool with a subset of commands (--contents, --control, --extract, --info)
* Create package from debpkg.yml specfile (like packager.io without cruft)
* GPG sign package
* GPG verify package`

	// Empty
	deb := New()
	defer deb.Close()

	deb.SetName("debpkg")
	deb.SetVersion("0.0.0")
	deb.SetMaintainer("Jerry Jacobs")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://github.com/xor-gate/debpkg")
	deb.SetShortDescription("Golang package for creating (gpg signed) debian packages")
	deb.SetDescription(controlDescr)
	deb.SetArchitecture("amd64")

	assert.Equal(t, controlExpect, deb.control.String(0))
}

// Test correct output of a control file Installed-Size property
func TestControlInstalledSize(t *testing.T) {
	// Empty
	deb := New()
	defer deb.Close()

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")

	// 1KByte
	controlExpect1K := `Package: 
Version: 0.0.0
Architecture: amd64
Maintainer:  <>
Installed-Size: 1
Description: 
`
	assert.Equal(t, controlExpect1K, deb.control.String(1024))

	// 2Kbyte
	controlExpect2K := `Package: 
Version: 0.0.0
Architecture: amd64
Maintainer:  <>
Installed-Size: 2
Description: 
`
	assert.Equal(t, controlExpect2K, deb.control.String(1025))
	assert.Equal(t, controlExpect2K, deb.control.String(2048))
}

func TestControlFileExtraString(t *testing.T) {
	deb := New()
	defer deb.Close()

	deb.SetName("debpkg-control-file-extra-string")
	deb.SetArchitecture("all")
	deb.SetDescription("bla bla\n")

	// BUG SetDescription should add the newline itself, and must not be left empty
	// dpkg: error processing archive /tmp/TestControlFileExtra.deb (--install):
	// parsing file '/var/lib/dpkg/tmp.ci/control' near line 6 package 'debpkg-control-file-extra:any':
	// end of file during value of field 'Description' (missing final newline)

	deb.AddControlExtraString("preinst", `#!/bin/sh
	echo "preinst: hello world from debpkg!"`)
	deb.AddControlExtraString("postinst", `#!/bin/sh
	echo "postinst: hello world from debpkg!"`)
	deb.AddControlExtraString("prerm", `#!/bin/sh
	echo "prerm: hello world from debpkg!"`)
	deb.AddControlExtraString("postrm", `#!/bin/sh
	echo "postrm: hello world from debpkg!"`)

	assert.Nil(t, testWrite(t, deb))
}

func TestControlFileExtra(t *testing.T) {
	deb := New()
	defer deb.Close()

	const script = `#!/bin/sh
echo "hello world from debpkg"
`

	filepath, err := test.WriteTempFile(t.Name()+".sh", script)
	assert.Nil(t, err)

	deb.SetName("debpkg-control-file-extra")
	deb.SetArchitecture("all")
	deb.SetDescription("bla bla\n")

	// BUG SetDescription should add the newline itself, and must not be left empty
	// dpkg: error processing archive /tmp/TestControlFileExtra.deb (--install):
	// parsing file '/var/lib/dpkg/tmp.ci/control' near line 6 package 'debpkg-control-file-extra:any':
	// end of file during value of field 'Description' (missing final newline)

	deb.AddControlExtra("preinst", filepath)
	deb.AddControlExtra("postinst", filepath)
	deb.AddControlExtra("prerm", filepath)
	deb.AddControlExtra("postrm", filepath)

	// FIXME AddControlExtra seems to add the full TMPDIR directory inside the control.tar.gz
	// E.g on mac var/folders/s5/x8wc0jqd387_sg4py6tg6xq00000gn/T/debpkg-test921120249

	assert.Nil(t, testWrite(t, deb))
}
