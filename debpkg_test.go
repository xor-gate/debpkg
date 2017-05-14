package debpkg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"golang.org/x/crypto/openpgp"
	"testing"
)

var e *openpgp.Entity

func init() {
	// Create random new GPG identity for signage
	e, _ = openpgp.NewEntity("Foo Bar", "", "foo@bar.com", nil)
}

// TestConfig verifies the specfile is correctly loaded
func TestConfig(t *testing.T) {
	deb := New()

	err := deb.Config("debpkg.yml")
	if err != nil {
		t.Errorf("Unable to open debpkg.yml in CWD: %v", err)
		return
	}

	if deb.control.info.version.full != "7.6.5" {
		t.Errorf("Unexpected deb.control.info.version.full: %s", deb.control.info.version.full)
		return
	}

	if deb.control.info.maintainer != "Foo Bar" {
		t.Errorf("Unexpected deb.control.info.maintainer: %s", deb.control.info.maintainer)
		return
	}

	if deb.control.info.maintainerEmail != "foo@bar.com" {
		t.Errorf("Unexpected deb.control.info.maintainerEmail: %s", deb.control.info.maintainerEmail)
		return
	}

	if deb.control.info.homepage != "https://github.com/xor-gate/debpkg" {
		t.Errorf("Unexpected deb.control.info.homepage: %s", deb.control.info.homepage)
		return
	}

	if deb.control.info.descrShort != "This is a short description" {
		t.Error("Unexpected short description")
		return
	}
}

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

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := deb.control.String(0)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
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

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/debpkg.git")
	deb.SetVcsBrowser("https://github.com/xor-gate/debpkg")
	control := deb.control.String(0)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
}

// Test correct output of the control file when SetVersion* functions are called
// Only the mandatory fields are exported then, this behaviour is checked
func TestControlFileSetVersionMajorMinorPatch(t *testing.T) {
	// Empty
	deb := New()

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
	control := deb.control.String(0)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}

	// Set full version string, this will overwrite the set SetVersion{Major,Minor,Patch} string
	deb.SetVersion("7.8.9")
	control = deb.control.String(0)

	controlExpectFullVersion := `Package: foobar
Version: 7.8.9
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Description: 
`

	if control != controlExpectFullVersion {
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

	deb.SetName("debpkg")
	deb.SetVersion("0.0.0")
	deb.SetMaintainer("Jerry Jacobs")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("https://github.com/xor-gate/debpkg")
	deb.SetShortDescription("Golang package for creating (gpg signed) debian packages")
	deb.SetDescription(controlDescr)
	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := deb.control.String(0)

	if control != controlExpect {
		t.Error("Unexpected control file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(controlExpect), controlExpect, len(control), control)
	}
}

// Test creation of empty digest
func TestDigestCreateEmpty(t *testing.T) {
	// FIXME it seems whe digesting the data buf the whole tarball will go corrupt...
	/*
	   	digestExpect := `Version: 4
	   Signer:
	   Date:
	   Role: builder
	   Files:
	   	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	   	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 control.tar.gz
	   	d41d8cd98f00b204e9800998ecf8427e da39a3ee5e6b4b0d3255bfef95601890afd80709 0 data.tar.gz
	   `
	*/
	digestExpect := `Version: 4
Signer: 
Date: 
Role: builder
Files: 
	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	0 0 0 control.tar.gz
	0 0 0 data.tar.gz
`

	deb := New()
	digest := createDigestFileString(deb)

	if digest != digestExpect {
		t.Error("Unexpected digest file")
		fmt.Printf("--- expected (len %d):\n'%s'\n--- got (len %d):\n'%s'---\n", len(digestExpect), digestExpect, len(digest), digest)
	}
}

// TestDirectory verifies adding a single directory recursive to the package
func TestAddDirectory(t *testing.T) {
	deb := New()
	deb.SetName("debpkg-test-add-directory")
	deb.SetArchitecture("all")
	err := deb.AddDirectory("vendor")
	if err != nil {
		t.Errorf("Error adding directory '.': %v", err)
		return
	}

	err = deb.Write("debpkg-test-add-directory.deb")
	if err != nil {
		t.Errorf("Error writing debfile: %v", err)
		return
	}
}

func TestWriteSignedEmpty(t *testing.T) {
	deb := New()

	// WriteSigned package
	err := deb.WriteSigned("debpkg-test-signed-empty.deb", e, "00000000")
	if err != nil {
		t.Errorf("Error in writing signed package: %v", err)
	}
}

func TestWrite(t *testing.T) {
	deb := New()

	deb.SetName("debpkg-test")
	deb.SetArchitecture("all")
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

	err := deb.Write("debpkg-test.deb")
	if err != nil {
		t.Errorf("Error in writing unsigned package: %v", err)
	}
}

func TestWriteSigned(t *testing.T) {
	deb := New()

	deb.SetName("debpkg-test-signed")
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
	deb.SetPriority(PriorityRequired)
	deb.SetConflicts("bash")
	deb.SetProvides("boembats")

	deb.AddFile("debpkg.go")

	// WriteSigned the package
	err := deb.WriteSigned("debpkg-test-signed.deb", e, "00000000")
	if err != nil {
		t.Errorf("Error in writing unsigned package: %v", err)
	}
}

func TestWriteError(t *testing.T) {
	deb := New()
	err := deb.Write("")
	if err == nil {
		t.Errorf("deb.Write shouldnt return nil")
	}
	deb.control.info.name = "pkg"
	if err := deb.Write(""); err == nil {
		t.Errorf("deb.Write shouldnt return nil")
	}
}

func ExampleWrite() {
	deb := New()

	deb.SetName("foobar")
	deb.SetVersion("1.2.3")
	deb.SetArchitecture("amd64")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepage("http://foobar.com")

	deb.SetShortDescription("Minimal foo bar package")
	deb.SetDescription("Foo bar package doesn't do anything")

	deb.AddFile("debpkg.go")
	fmt.Println(deb.Write("foobar.deb"))

	// Output: <nil>
}

func dpkg(cmd, action, filename string) error {
	args := []string{"--"+action, filename}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return err
	}
	return nil
}

func TestReadWithNativeDpkg(t *testing.T) {
	dpkgCmd, err := exec.LookPath("dpkg")
	if err != nil || dpkgCmd == "" {
		fmt.Println("Skip test, unable to find dpkg in PATH")
		return
	}

	debs, err := filepath.Glob("*.deb")
	for _, deb := range debs {
		err = dpkg(dpkgCmd, "info", deb)
		if err != nil {
			t.Errorf("dpkg --info failed on " + deb)
		}
		fmt.Println("dpkg --info passed on " + deb)

		dpkg(dpkgCmd, "contents", deb)
		if err != nil {
			t.Errorf("dpkg --contents failed on " + deb)
		}
		fmt.Println("dpkg --contents passed on " + deb)
	}
}
