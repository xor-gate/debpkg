// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xor-gate/debpkg/internal/test"
)

// TestExampleConfig verifies if the config example in the root is correctly loaded
func TestExampleConfig(t *testing.T) {
	const configFile = `name: debpkg
version: 7.6.5
architecture: all
maintainer: Deb Pkg
maintainer_email: deb@pkg.com
homepage: https://github.com/xor-gate/debpkg
section: devel
priority: standard
depends: lsb-release
recommends: nano
suggests: curl
conflicts: pico
provides: editor
replaces: vim
built_using: golang
description:
   short: This is a short description
   long: >
       Bla bla
       Bla Bla
       .
       Dusse
files:
  - file: LICENSE
    dest: {{.DATAROOTDIR}}/foobar/LICENSE
  - file: debpkg.go
  - file: debpkg_test.go
  - file: README.md
    dest: {{.DATAROOTDIR}}/foobar/README.md
  - dest: /bin/hello
    content: >
      #!/bin/bash
      echo "hello"
directories:
  - ./internal
emptydirs:
  - /var/cache/foobar
control_extra:
  postrm: Makefile
  prerm: Makefile
  postinst: Makefile
  preinst: Makefile
`
	filepath, err := test.WriteTempFile("debpkg.yml", configFile)
	assert.Nil(t, err)

	deb := New()
	defer deb.Close()

	assert.Nil(t, deb.Config(filepath))
	assert.Equal(t, "7.6.5", deb.Version,
		"Unexpected deb.Version")
	assert.Equal(t, "7.6.5", deb.control.Info.Version,
		"Unexpected deb.control.Info.version")
	assert.Equal(t, "Deb Pkg", deb.control.Info.Maintainer,
		"Unexpected deb.control.Info.maintainer")
	assert.Equal(t, "deb@pkg.com", deb.control.Info.MaintainerEmail,
		"Unexpected deb.control.Info.maintainerEmail")
	assert.Equal(t, "https://github.com/xor-gate/debpkg", deb.control.Info.Homepage,
		"Unexpected deb.control.Info.homepage")
	assert.Equal(t, "This is a short description", deb.control.Info.DescrShort,
		"Unexpected short description")
	assert.Equal(t, "golang", deb.control.Info.BuiltUsing,
		"unexpected built using")
	assert.Equal(t, "devel", deb.control.Info.Section,
		"unexpected section")
	assert.Equal(t, PriorityStandard, deb.control.Info.Priority,
		"unexpected priority")

	assert.Nil(t, testWrite(t, deb))
}

func TestExampleConfigWithControlExtraContent(t *testing.T) {
	const configFile = `name: foo-bar
version: 1.2.3
architecture: amd64
maintainer: Mr. Foo Bar
maintainer_email: foo@bar.org
homepage: https://www.debian.org
section: net
priority: important
control_extra:
  postrm: >
    #!/bin/bash
    echo "post rm!!"
  prerm: >
    #!/bin/bash
    echo "pre rm!!"
  postinst: >
    #!/bin/bash
    echo "post inst!!"
  preinst: 	>
    #!/bin/bash
    echo "pre inst!!"
`
	filepath, err := test.WriteTempFile(t.Name()+".yml", configFile)
	assert.Nil(t, err)

	deb := New()
	defer deb.Close()

	assert.Nil(t, deb.Config(filepath))
	assert.Equal(t, "1.2.3", deb.Version,
		"Unexpected deb.Version")
	assert.Equal(t, "1.2.3", deb.control.Info.Version,
		"Unexpected deb.control.Info.version")
	assert.Equal(t, "Mr. Foo Bar", deb.control.Info.Maintainer,
		"Unexpected deb.control.Info.maintainer")
	assert.Equal(t, "foo@bar.org", deb.control.Info.MaintainerEmail,
		"Unexpected deb.control.Info.maintainerEmail")
	assert.Equal(t, "https://www.debian.org", deb.control.Info.Homepage,
		"Unexpected deb.control.Info.homepage")
	assert.Equal(t, "net", deb.control.Info.Section,
		"unexpected section")
	assert.Equal(t, PriorityImportant, deb.control.Info.Priority,
		"unexpected priority")

	assert.Nil(t, testWrite(t, deb))
}

func TestExampleConfigWithConfigFile(t *testing.T) {
	const configFile = `name: bar-bar
version: 1.1.1
architecture: amd64
maintainer: Mr. Foo Bar
maintainer_email: foo@bar.org
homepage: https://www.debian.org
section: net
priority: important
files:
  - dest: /etc/hello
    conffile: true
    content: >
      #!/bin/bash
      echo "hello"
  - dest: /my/awesome/makefile
    conffile: true
    file: Makefile

`
	filepath, err := test.WriteTempFile(t.Name()+".yml", configFile)
	assert.Nil(t, err)

	deb := New()
	defer deb.Close()

	assert.Nil(t, deb.Config(filepath))
	assert.Equal(t, "1.1.1", deb.Version,
		"Unexpected deb.Version")
	assert.Equal(t, "1.1.1", deb.control.Info.Version,
		"Unexpected deb.control.Info.version")
	assert.Equal(t, "/etc/hello\n/my/awesome/makefile\n", deb.control.conffiles)

	assert.Nil(t, testWrite(t, deb))
}

func TestDefaultConfig(t *testing.T) {
	filepath, err := test.WriteTempFile(t.Name()+".yml", "")
	assert.Nil(t, err)

	deb := New()
	defer deb.Close()

	assert.Nil(t, deb.Config(filepath))

	assert.Equal(t, "any", deb.control.Info.Architecture,
		"unexpected architecture")
	assert.Equal(t, "anonymous", deb.control.Info.Maintainer,
		"unexpected maintainer")
	assert.Equal(t, "anon@foo.bar", deb.control.Info.MaintainerEmail,
		"unexpected maintainer email")
	assert.Equal(t, "https://www.google.com", deb.control.Info.Homepage,
		"unexpected homepage")
	assert.Equal(t, PriorityOptional, deb.control.Info.Priority,
		"unexpected priority")
	assert.Equal(t, "0.1.0+dev", deb.Version,
		"unexpected version")
	assert.Equal(t, "0.1.0+dev", deb.control.Info.Version,
		"Unexpected deb.control.Info.version")
	assert.Equal(t, "misc", deb.control.Info.Section,
		"unexpected section")
	assert.Equal(t, "unknown", deb.control.Info.Name,
		"unexpected name")
	assert.Equal(t, runtime.Version(), deb.control.Info.BuiltUsing,
		"unexpected built using")
	assert.Equal(t, "-", deb.control.Info.DescrShort,
		"unexpected short description")
	assert.Equal(t, " -", deb.control.Info.Descr,
		"unexpected long description")
}

func TestNonExistingConfig(t *testing.T) {
	deb := New()
	defer deb.Close()

	assert.NotNil(t, deb.Config("/non/existent/config/file"))
}

func TestInvalidYAML(t *testing.T) {
	deb := New()
	defer deb.Close()

	const configFile = `name: debpkg
	foo: bar
	`
	filepath, err := test.WriteTempFile(t.Name()+".yml", configFile)
	assert.Nil(t, err)
	assert.NotNil(t, deb.Config(filepath))
}

func TestInvalidTemplateVar(t *testing.T) {
	deb := New()
	defer deb.Close()

	const configFile = `name: debpkg
foo: {{.bar}}
	`
	filepath, err := test.WriteTempFile(t.Name()+".yml", configFile)
	assert.Nil(t, err)
	assert.NotNil(t, deb.Config(filepath))
}
