package debpkg

import (
	"io/ioutil"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExampleConfig verifies if the config example in the root is correctly loaded
func TestExampleConfig(t *testing.T) {
	deb := New()

	err := deb.Config("debpkg.yml")
	if err != nil {
		t.Errorf("Unable to open debpkg.yml in CWD: %v", err)
	}
	assert.Equal(t, "7.6.5", deb.control.info.version.full,
		"Unexpected deb.control.info.version.full")
	assert.Equal(t, "Foo Bar", deb.control.info.maintainer,
		"Unexpected deb.control.info.maintainer")
	assert.Equal(t, "foo@bar.com", deb.control.info.maintainerEmail,
		"Unexpected deb.control.info.maintainerEmail")
	assert.Equal(t, "https://github.com/xor-gate/debpkg", deb.control.info.homepage,
		"Unexpected deb.control.info.homepage")
	assert.Equal(t, "This is a short description", deb.control.info.descrShort,
		"Unexpected short description")
	assert.Equal(t, "golang", deb.control.info.builtUsing,
		"unexpected built using")
	assert.Equal(t, "devel", deb.control.info.section,
		"unexpected section")
	assert.Equal(t, PriorityStandard, deb.control.info.priority,
		"unexpected priority")
}

func TestDefaultConfig(t *testing.T) {
	f, err := ioutil.TempFile("", "config")
	if err != nil {
		t.Errorf("unexpected error creating tempfile: %v", err)
	}
	f.Close()
	deb := New()
	if err := deb.Config(f.Name()); err != nil {
		t.Errorf("Unexpected error during load of empty config: %v", err)
	}
	assert.Equal(t, "any", deb.control.info.architecture,
		"unexpected architecture")
	assert.Equal(t, "anonymous", deb.control.info.maintainer,
		"unexpected maintainer")
	assert.Equal(t, "anon@foo.bar", deb.control.info.maintainerEmail,
		"unexpected maintainer email")
	assert.Equal(t, "https://www.google.com", deb.control.info.homepage,
		"unexpected homepage")
	assert.Equal(t, PriorityOptional, deb.control.info.priority,
		"unexpected priority")
	assert.Equal(t, "0.1.0+dev", deb.control.info.version.full,
		"unexpected version")
	assert.Equal(t, "misc", deb.control.info.section,
		"unexpected section")
	assert.Equal(t, "unknown", deb.control.info.name,
		"unexpected name")
	assert.Equal(t, runtime.Version(), deb.control.info.builtUsing,
		"unexpected built using")
	assert.Equal(t, "-", deb.control.info.descrShort,
		"unexpected short description")
	assert.Equal(t, " -", deb.control.info.descr,
		"unexpected long description")

	/*
		deb.control.info.conflicts
		deb.control.info.provides
		deb.control.info.replaces
		deb.control.info.suggests
		vcs**/
}
