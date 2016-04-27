package debpkg_test

import (
	"fmt"
	"testing"
	"github.com/xor-gate/debpkg"
)

func TestConfig(t *testing.T) {
	deb := debpkg.New()

	err := deb.Config("debpkg.yml")
	fmt.Println(err)
}

func TestWrite(t *testing.T) {
	deb := debpkg.New()

	deb.SetName("test")
	deb.SetVersion("0.0.1")
	deb.SetMaintainer("Foo Bar")
	deb.SetMaintainerEmail("foo@bar.com")
	deb.SetHomepageUrl("https://foobar.com")
	deb.SetShortDescription("some awesome foobar pkg")
	deb.SetDescription("very very very very long description")

	deb.AddFile("debpkg.go")
	deb.AddDirectory("tests")

	deb.Sign()
	deb.Write("test.deb")
}
