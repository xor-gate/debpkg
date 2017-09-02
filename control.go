// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"

	"github.com/xor-gate/debpkg/internal/targzip"
)

type control struct {
	tgz       *targzip.TarGzip
	info      controlInfo
	extra     map[string][]byte // Extra files added to the control.tar.gz. Typical usage is for preinst, postinst, prerm, postrm.
	conffiles map[string][]byte // Conffiles which must be treated as configuration files
}

type controlInfoVersion struct {
	full  string // Full version string. E.g "0.1.2"
	major uint   // Major version number
	minor uint   // Minor version number
	patch uint   // Patch version number
}

type controlInfo struct {
	name            string
	version         controlInfoVersion
	architecture    string
	maintainer      string
	maintainerEmail string
	homepage        string
	depends         string
	recommends      string
	suggests        string
	conflicts       string
	provides        string
	replaces        string
	section         string
	priority        Priority
	descrShort      string  // Short package description
	descr           string  // Long package description
	vcsType         VcsType // E.g: "Svn", "Git" etcetera
	vcsURL          string  // E.g: git@github.com:xor-gate/debpkg.git
	vcsBrowser      string  // E.g: https://github.com/xor-gate/debpkg
	builtUsing      string  // E.g: gcc-4.6 (= 4.6.0-11)
}

// SetName sets the name of the binary package (mandatory)
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Package
func (deb *DebPkg) SetName(name string) {
	deb.control.info.name = name
}

// SetVersion sets the full version string (mandatory), or use SetVersion* functions for "major.minor.patch"
// The upstream_version may contain only alphanumerics ( A-Za-z0-9 ) and the characters . + - : ~
//  (full stop, plus, hyphen, colon, tilde) and should start with a digit.
// NOTE: When the full string is set the other SetVersion* function calls are ignored
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersion(version string) {
	deb.control.info.version.full = version
}

// SetVersionMajor sets the version major number
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionMajor(major uint) {
	deb.control.info.version.major = major
}

// SetVersionMinor sets the version minor number
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionMinor(minor uint) {
	deb.control.info.version.minor = minor
}

// SetVersionPatch sets the version patch level
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionPatch(patch uint) {
	deb.control.info.version.patch = patch
}

// SetArchitecture sets the architecture of the package where it can be installed.
//  E.g "i386, amd64, arm, any, all". See `dpkg-architecture -L` for all supported.
// Architecture: any
//    The generated binary package is an architecture dependent one usually in a compiled language.
// Architecture: all
//    The generated binary package is an architecture independent one usually consisting of text,
//    images, or scripts in an interpreted language.
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Architecture
// And: http://man7.org/linux/man-pages/man1/dpkg-architecture.1.html
func (deb *DebPkg) SetArchitecture(arch string) {
	deb.control.info.architecture = arch
}

// SetMaintainer (mandatory), sets the package maintainers name and surname. E.g: "Foo Bar"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (deb *DebPkg) SetMaintainer(maintainer string) {
	deb.control.info.maintainer = maintainer
}

// SetMaintainerEmail sets the package maintainers email address. E.g: "foo@bar.com"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (deb *DebPkg) SetMaintainerEmail(email string) {
	deb.control.info.maintainerEmail = email
}

// SetDepends sets the package dependencies. E.g: "lsb-release"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (deb *DebPkg) SetDepends(depends string) {
	deb.control.info.depends = depends
}

// SetRecommends sets the package recommendations. E.g: "aptitude"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (deb *DebPkg) SetRecommends(recommends string) {
	deb.control.info.recommends = recommends
}

// SetSuggests sets the package suggestions. E.g: "aptitude"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (deb *DebPkg) SetSuggests(suggests string) {
	deb.control.info.suggests = suggests
}

// SetConflicts sets one or more conflicting packages. E.g: "nano"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-conflicts
func (deb *DebPkg) SetConflicts(conflicts string) {
	deb.control.info.conflicts = conflicts
}

// SetProvides sets the type which the package provides. E.g: "editor"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-virtual
func (deb *DebPkg) SetProvides(provides string) {
	deb.control.info.provides = provides
}

// SetReplaces sets the names of packages which will be replaced. E.g: "pico"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html
func (deb *DebPkg) SetReplaces(replaces string) {
	deb.control.info.replaces = replaces
}

// SetPriority (recommended). Default set to debpkg.PriorityUnset
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Priority
// And: https://www.debian.org/doc/debian-policy/ch-archive.html#s-priorities
func (deb *DebPkg) SetPriority(priority Priority) {
	deb.control.info.priority = priority
}

// SetSection (recommended). E.g: editors
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Section
// And: https://www.debian.org/doc/debian-policy/ch-archive.html#s-subsections
func (deb *DebPkg) SetSection(section string) {
	deb.control.info.section = section
}

// SetHomepage sets the homepage URL of the package. E.g: "https://github.com/foo/bar"
func (deb *DebPkg) SetHomepage(url string) {
	deb.control.info.homepage = url
}

// SetShortDescription sets the single line synopsis. E.g: "My awesome foo bar baz tool"
func (deb *DebPkg) SetShortDescription(descr string) {
	deb.control.info.descrShort = descr
}

// SetDescription sets the extended description over several lines. E.g:
// "This tool will calculation the most efficient way to world domination"
// NOTE: The debian control file has a special formatting of the long description
//        this function replaces newlines with a newline and a space.
func (deb *DebPkg) SetDescription(descr string) {
	deb.control.info.descr = " " + strings.Replace(descr, "\n", "\n ", -1)
}

// SetVcsType sets the version control system (Vcs) type for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsType(vcs VcsType) {
	deb.control.info.vcsType = vcs
}

// SetVcsURL sets the version control system (Vcs) URL for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsURL(url string) {
	deb.control.info.vcsURL = url
}

// SetVcsBrowser sets the version control system (Vcs) browsable source-tree URL for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsBrowser(url string) {
	deb.control.info.vcsBrowser = url
}

// SetBuiltUsing incorporate parts of other packages when built but do not have to depend on those packages.
// A package using the source code from the gcc-4.6-source binary package built from the gcc-4.6 source package
// would have this field in its control file:
//  Built-Using: gcc-4.6 (= 4.6.0-11)
// A package including binaries from grub2 and loadlin would have this field in its control file:
//  Built-Using: grub2 (= 1.99-9), loadlin (= 1.6e-1)
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-built-using
func (deb *DebPkg) SetBuiltUsing(info string) {
	deb.control.info.builtUsing = info
}

// AddControlExtraString is the same as AddControlExtra except it uses a string input
func (deb *DebPkg) AddControlExtraString(name, s string) error {
	s = strings.Replace(s, "\r\n", "\n", -1)
	return deb.control.tgz.AddFileFromBuffer(name, []byte(s))
}

// AddControlExtra allows the advanced user to add custom script to the control.tar.gz Typical usage is
//  for preinst, postinst, postrm, prerm: https://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
// And: https://www.debian.org/doc/manuals/maint-guide/dother.en.html#maintscripts
func (deb *DebPkg) AddControlExtra(name, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return deb.AddControlExtraString(name, string(b))
}

// AddConffile adds a file to the conffiles so it is treated as configuration files. Configuration files are not
// overwritten during an update unless specified.
func (deb *DebPkg) AddConffile(filename string, r io.Reader) {
}

// verify the control file for validity
func (c *control) verify() error {
	if c.info.name == "" {
		return fmt.Errorf("empty package name")
	}
	if c.info.architecture == "" {
		return fmt.Errorf("empty architecture")
	}
	return nil
}

func createControlTarGz(deb *DebPkg) error {
	controlFile := []byte(deb.control.String(deb.data.tgz.Written()))
	if err := deb.control.tgz.AddFileFromBuffer("control", controlFile); err != nil {
		return err
	}
	if err := deb.control.tgz.AddFileFromBuffer("md5sums", []byte(deb.data.md5sums)); err != nil {
		return err
	}
	return nil
}

// Generate version string (e.g "1.2.3") from major,minor patch or from full version
func (c *control) version() string {
	if c.info.version.full != "" {
		return c.info.version.full
	}
	return fmt.Sprintf("%d.%d.%d",
		c.info.version.major,
		c.info.version.minor,
		c.info.version.patch)
}

func (c *control) size() int64 {
	return c.tgz.Size()
}

// Create control file for control.tar.gz
func (c *control) String(installedSize uint64) string {
	var o string

	o += fmt.Sprintf("Package: %s\n", c.info.name)
	o += fmt.Sprintf("Version: %s\n", c.version())
	o += fmt.Sprintf("Architecture: %s\n", c.info.architecture)
	o += fmt.Sprintf("Maintainer: %s <%s>\n",
		c.info.maintainer,
		c.info.maintainerEmail)
	o += fmt.Sprintf("Installed-Size: %d\n", uint64(math.Ceil(float64(installedSize)/1024)))

	if c.info.section != "" {
		o += fmt.Sprintf("Section: %s\n", c.info.section)
	}
	if c.info.priority != PriorityUnset {
		o += fmt.Sprintf("Priority: %s\n", c.info.priority)
	}
	if c.info.homepage != "" {
		o += fmt.Sprintf("Homepage: %s\n", c.info.homepage)
	}
	if c.info.vcsType != VcsTypeUnset && c.info.vcsURL != "" {
		o += fmt.Sprintf("Vcs-%s: %s\n", c.info.vcsType, c.info.vcsURL)
	}
	if c.info.vcsBrowser != "" {
		o += fmt.Sprintf("Vcs-Browser: %s\n", c.info.vcsBrowser)
	}
	if c.info.builtUsing != "" {
		o += fmt.Sprintf("Built-Using: %s\n", c.info.builtUsing)
	}

	if c.info.depends != "" {
		o += fmt.Sprintf("Depends: %s\n", c.info.depends)
	}
	if c.info.recommends != "" {
		o += fmt.Sprintf("Recommends: %s\n", c.info.recommends)
	}
	if c.info.suggests != "" {
		o += fmt.Sprintf("Suggests: %s\n", c.info.suggests)
	}
	if c.info.conflicts != "" {
		o += fmt.Sprintf("Conflicts: %s\n", c.info.conflicts)
	}
	if c.info.provides != "" {
		o += fmt.Sprintf("Provides: %s\n", c.info.provides)
	}
	if c.info.replaces != "" {
		o += fmt.Sprintf("Replaces: %s\n", c.info.replaces)
	}

	o += fmt.Sprintf("Description: %s\n", c.info.descrShort)
	o += fmt.Sprintf("%s", c.info.descr)

	return o
}
