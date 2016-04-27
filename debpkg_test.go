package debpkg

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"testing"
)

func TestConfig(t *testing.T) {
	deb := New()

	err := deb.Config("yml")
	fmt.Println(err)
}

// Test correct output of a empty control file when no DepPkg Set* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileEmpty(t *testing.T) {
	controlExpect := `Package: 
Version: 0.0.0
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Homepage: 
Description: 
`
	// Empty
	deb := New()

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := createControlFileString(deb)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
}

// Test correct output of the control file when SetVersion* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileSetVersionMajorMinorPatch(t *testing.T) {
	controlExpect := `Package: 
Version: 1.2.3
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Homepage: 
Description: 
`
	// Empty
	deb := New()

	deb.SetVersionMajor(1)
	deb.SetVersionMinor(2)
	deb.SetVersionPatch(3)

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := createControlFileString(deb)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
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
 * GPG verify package
`

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

	deb.SetName("debpkg")
	deb.SetVersion("0.0.0")
	deb.SetMaintainer("Jerry Jacobs")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://github.com/xor-gate/debpkg")
	deb.SetShortDescription("Golang package for creating (gpg signed) debian packages")
	deb.SetDescription(controlDescr)
	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := createControlFileString(deb)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
}

func TestSign(t *testing.T) {
	deb := New()

	// Create random new GPG identity for signage
	var e *openpgp.Entity
	e, err := openpgp.NewEntity("Foo Bar", "", "foo@bar.com", nil)
	if err != nil {
		t.Error("Unexpected New GPG Entity")
		return
	}

	// Sign
	deb.Sign(e, "00000000")
}

func TestWrite(t *testing.T) {
	deb := New()

	deb.SetName("test")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://foobar.com")
	deb.SetShortDescription("some awesome foobar pkg")
	deb.SetDescription("very very very very long description")

	// Set version control system info for control file
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/secdl")
	deb.SetVcsBrowser("https://github.com/xor-gate/secdl")

	deb.AddFile("debpkg.go")
	// FIXME deb.AddDirectory("tests")

	deb.Write("debpkg.deb")
}
