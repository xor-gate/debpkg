package debpkg

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	deb := New()

	err := deb.Config("yml")
	fmt.Println(err)
}

// Test correct output of a empty control file when no DepPkg Set* functions are called
func TestControlFileEmpty(t *testing.T) {
controlExpect := `Package: 
Architecture: amd64
Maintainer:  <>
Installed-Size: 0
Section: 
Priority: 
Homepage: 
Description: 
`
	// Empty
	deb := New()

	// architecture is auto-set when empty, this makes sure it is always set to amd64
	deb.SetArchitecture("amd64")
	control := createControlFileString(deb)

	if control != controlExpect {
		t.Error("Unexpected empty control file")
	}
}

func TestWrite(t *testing.T) {
	deb := New()

	deb.SetName("test")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepageUrl("https://foobar.com")
	deb.SetShortDescription("some awesome foobar pkg")
	deb.SetDescription("very very very very long description")

	// Set version control system info for control file
	deb.SetVcsType(VcsTypeGit)
	deb.SetVcsURL("https://github.com/xor-gate/secdl")
	deb.SetVcsBrowser("https://github.com/xor-gate/secdl")

	deb.AddFile("go")
	deb.AddDirectory("tests")

	deb.Sign()
	deb.Write("deb")
}
