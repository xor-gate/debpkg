// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"math"
	"strings"

	"github.com/xor-gate/debpkg/internal/targzip"
)

// ControlFile contains the Debian package metadata and information
type ControlFile struct {
	Info               ControlFileInfo
	tgz                *targzip.TarGzip
	conffiles          string // List of configuration-files
	hasCustomConffiles bool
}

// ControlFileInfo contains the information of a Debian package
type ControlFileInfo struct {
	Name            string
	Version         string
	Architecture    string
	Maintainer      string
	MaintainerEmail string
	Homepage        string
	Depends         string
	Recommends      string
	Suggests        string
	Conflicts       string
	Provides        string
	Replaces        string
	Section         string
	Priority        Priority
	DescrShort      string  // Short package description
	Descr           string  // Long package description
	VcsType         VcsType // E.g: "Svn", "Git" etcetera
	VcsURL          string  // E.g: git@github.com:xor-gate/debpkg.git
	VcsBrowser      string  // E.g: https://github.com/xor-gate/debpkg
	BuiltUsing      string  // E.g: gcc-4.6 (= 4.6.0-11)
}

// SetName sets the name of the binary package (mandatory)
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Package
func (pkg *Package) SetName(name string) {
	pkg.Name = name
	// TODO control file is not able to access pkg.Name so we copy it for now
	pkg.control.Info.Name = name
}

// SetVersion sets the full version string (mandatory), or use SetVersion* functions for "major.minor.patch"
// The upstream_version may contain only alphanumerics ( A-Za-z0-9 ) and the characters . + - : ~
//  (full stop, plus, hyphen, colon, tilde) and should start with a digit.
// NOTE: When the full string is set the other SetVersion* function calls are ignored
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (pkg *Package) SetVersion(version string) {
	pkg.Version = version
	// TODO control.Info is not able to access pkg.Version so we copy it for now
	pkg.control.Info.Version = version
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
func (pkg *Package) SetArchitecture(arch string) {
	pkg.control.Info.Architecture = arch
}

// SetMaintainer (mandatory), sets the package maintainers name and surname. E.g: "Foo Bar"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (pkg *Package) SetMaintainer(maintainer string) {
	pkg.control.Info.Maintainer = maintainer
}

// SetMaintainerEmail sets the package maintainers email address. E.g: "foo@bar.com"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (pkg *Package) SetMaintainerEmail(email string) {
	pkg.control.Info.MaintainerEmail = email
}

// SetDepends sets the package dependencies. E.g: "lsb-release"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (pkg *Package) SetDepends(depends string) {
	pkg.control.Info.Depends = depends
}

// SetRecommends sets the package recommendations. E.g: "aptitude"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (pkg *Package) SetRecommends(recommends string) {
	pkg.control.Info.Recommends = recommends
}

// SetSuggests sets the package suggestions. E.g: "aptitude"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-binarydeps
func (pkg *Package) SetSuggests(suggests string) {
	pkg.control.Info.Suggests = suggests
}

// SetConflicts sets one or more conflicting packages. E.g: "nano"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-conflicts
func (pkg *Package) SetConflicts(conflicts string) {
	pkg.control.Info.Conflicts = conflicts
}

// SetProvides sets the type which the package provides. E.g: "editor"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-virtual
func (pkg *Package) SetProvides(provides string) {
	pkg.control.Info.Provides = provides
}

// SetReplaces sets the names of packages which will be replaced. E.g: "pico"
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html
func (pkg *Package) SetReplaces(replaces string) {
	pkg.control.Info.Replaces = replaces
}

// SetPriority (recommended). Default set to debpkg.PriorityUnset
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Priority
// And: https://www.debian.org/doc/debian-policy/ch-archive.html#s-priorities
func (pkg *Package) SetPriority(priority Priority) {
	pkg.control.Info.Priority = priority
}

// SetSection (recommended). E.g: editors
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Section
// And: https://www.debian.org/doc/debian-policy/ch-archive.html#s-subsections
func (pkg *Package) SetSection(section string) {
	pkg.control.Info.Section = section
}

// SetHomepage sets the homepage URL of the package. E.g: "https://github.com/foo/bar"
func (pkg *Package) SetHomepage(url string) {
	pkg.control.Info.Homepage = url
}

// SetShortDescription sets the single line synopsis. E.g: "My awesome foo bar baz tool"
func (pkg *Package) SetShortDescription(descr string) {
	pkg.control.Info.DescrShort = descr
}

// SetDescription sets the extended description over several lines. E.g:
// "Debpkg calculates the most efficient way to world domination"
// NOTE: The debian control file has a special formatting of the long description
//        this function replaces newlines with a newline and a space.
func (pkg *Package) SetDescription(descr string) {
	pkg.control.Info.Descr = " " + strings.Replace(descr, "\n", "\n ", -1)
}

// SetVcsType sets the version control system (Vcs) type for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (pkg *Package) SetVcsType(vcs VcsType) {
	pkg.control.Info.VcsType = vcs
}

// SetVcsURL sets the version control system (Vcs) URL for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (pkg *Package) SetVcsURL(url string) {
	pkg.control.Info.VcsURL = url
}

// SetVcsBrowser sets the version control system (Vcs) browsable source-tree URL for the source package.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (pkg *Package) SetVcsBrowser(url string) {
	pkg.control.Info.VcsBrowser = url
}

// SetBuiltUsing incorporate parts of other packages when built but do not have to depend on those packages.
// A package using the source code from the gcc-4.6-source binary package built from the gcc-4.6 source package
// would have this field in its control file:
//  Built-Using: gcc-4.6 (= 4.6.0-11)
// A package including binaries from grub2 and loadlin would have this field in its control file:
//  Built-Using: grub2 (= 1.99-9), loadlin (= 1.6e-1)
// See: https://www.debian.org/doc/debian-policy/ch-relationships.html#s-built-using
func (pkg *Package) SetBuiltUsing(info string) {
	pkg.control.Info.BuiltUsing = info
}

// Verify the control file for validity
func (c *ControlFile) Verify() error {
	if c.Info.Name == "" {
		return fmt.Errorf("empty package name")
	}
	if c.Info.Architecture == "" {
		return fmt.Errorf("empty architecture")
	}
	if c.Info.Version == "" {
		return fmt.Errorf("empty package version")
	}
	return nil
}

// MarkConfigFile adds the dest absolute filename to conffiles
func (c *ControlFile) MarkConfigFile(dest string) error {
	if dest == "" {
		return fmt.Errorf("config file cannot be empty")
	}
	c.conffiles += dest + "\n"
	return nil
}

// finalizeControlFile creates the actual control-file, adds MD5-sums and stores
// conffiles
func (c *ControlFile) finalizeControlFile(d *data) error {
	if !c.hasCustomConffiles {
		if err := c.tgz.AddFileFromBuffer("conffiles", []byte(c.conffiles)); err != nil {
			return err
		}
	}
	controlFile := []byte(c.String(d.tgz.Written()))
	if err := c.tgz.AddFileFromBuffer("control", controlFile); err != nil {
		return err
	}
	if err := c.tgz.AddFileFromBuffer("md5sums", []byte(d.md5sums)); err != nil {
		return err
	}
	return nil
}

// Create control file for control.tar.gz
func (c *ControlFile) String(installedSize uint64) string {
	var o string

	o += fmt.Sprintf("Package: %s\n", c.Info.Name)
	o += fmt.Sprintf("Version: %s\n", c.Info.Version)
	o += fmt.Sprintf("Architecture: %s\n", c.Info.Architecture)
	o += fmt.Sprintf("Maintainer: %s <%s>\n",
		c.Info.Maintainer,
		c.Info.MaintainerEmail)
	o += fmt.Sprintf("Installed-Size: %d\n", uint64(math.Ceil(float64(installedSize)/1024)))

	if c.Info.Section != "" {
		o += fmt.Sprintf("Section: %s\n", c.Info.Section)
	}
	if c.Info.Priority != PriorityUnset {
		o += fmt.Sprintf("Priority: %s\n", c.Info.Priority)
	}
	if c.Info.Homepage != "" {
		o += fmt.Sprintf("Homepage: %s\n", c.Info.Homepage)
	}
	if c.Info.VcsType != VcsTypeUnset && c.Info.VcsURL != "" {
		o += fmt.Sprintf("Vcs-%s: %s\n", c.Info.VcsType, c.Info.VcsURL)
	}
	if c.Info.VcsBrowser != "" {
		o += fmt.Sprintf("Vcs-Browser: %s\n", c.Info.VcsBrowser)
	}
	if c.Info.BuiltUsing != "" {
		o += fmt.Sprintf("Built-Using: %s\n", c.Info.BuiltUsing)
	}

	if c.Info.Depends != "" {
		o += fmt.Sprintf("Depends: %s\n", c.Info.Depends)
	}
	if c.Info.Recommends != "" {
		o += fmt.Sprintf("Recommends: %s\n", c.Info.Recommends)
	}
	if c.Info.Suggests != "" {
		o += fmt.Sprintf("Suggests: %s\n", c.Info.Suggests)
	}
	if c.Info.Conflicts != "" {
		o += fmt.Sprintf("Conflicts: %s\n", c.Info.Conflicts)
	}
	if c.Info.Provides != "" {
		o += fmt.Sprintf("Provides: %s\n", c.Info.Provides)
	}
	if c.Info.Replaces != "" {
		o += fmt.Sprintf("Replaces: %s\n", c.Info.Replaces)
	}

	o += fmt.Sprintf("Description: %s\n", c.Info.DescrShort)
	o += fmt.Sprintf("%s", c.Info.Descr)

	return o
}
