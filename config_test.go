package debpkg

import "testing"

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
