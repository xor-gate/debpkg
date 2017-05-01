// Copyright 2016 Jerry Jacobs. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"go/build"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blakesmith/ar"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
)

// Priority for Debian package
type Priority string

// Package Priority
const (
	PriorityUnset     Priority = ""          // Priority field is skipped
	PriorityRequired  Priority = "required"  // Packages which are necessary for the proper functioning of the system
	PriorityImportant Priority = "important" // Important programs, including those which one would expect to find on any Unix-like system
	PriorityStandard  Priority = "standard"  // These packages provide a reasonably small but not too limited character-mode system
	PriorityOptional  Priority = "optional"  // This is all the software that you might reasonably want to install if you didn't know what it was and don't have specialized requirements
)

// VcsType for Debian package supported version control system (Vcs) types
type VcsType string

// Package VcsType
const (
	VcsTypeUnset      VcsType = ""      // VcsType field is skipped
	VcsTypeArch       VcsType = "Arch"  // Arch
	VcsTypeBazaar     VcsType = "Bzr"   // Bazaar
	VcsTypeDarcs      VcsType = "Darcs" // Darcs
	VcsTypeGit        VcsType = "Git"   // Git
	VcsTypeMercurial  VcsType = "Hg"    // Mercurial
	VcsTypeMonotone   VcsType = "Mtn"   // Monotone
	VcsTypeSubversion VcsType = "Svn"   // Subversion
)

const debPkgDebianBinary = "2.0\n"
const debPkgDigestDefaultHash = crypto.SHA1
const debPkgDigestVersion = 4
const debPkgDigestRole = "builder"

type debPkgData struct {
	size    int64
	md5sums string
	buf     *bytes.Buffer
	tw      *tar.Writer
	gw      *gzip.Writer
}

type debPkgVersion struct {
	full  string // Full version string. E.g "0.1.2"
	major uint   // Major version number
	minor uint   // Minor version number
	patch uint   // Patch version number
}

type debPkgControlInfo struct {
	name            string
	version         debPkgVersion
	architecture    string
	maintainer      string
	maintainerEmail string
	homepage        string
	suggests        string
	conflicts       string
	replaces        string
	provides        string
	section         string
	priority        Priority
	descrShort      string  // Short package description
	descr           string  // Long package description
	vcsType         VcsType // E.g: "Svn", "Git" etcetera
	vcsURL          string  // E.g: git@github.com:xor-gate/debpkg.git
	vcsBrowser      string  // E.g: https://github.com/xor-gate/debpkg
	builtUsing      string  // E.g: gcc-4.6 (= 4.6.0-11)
}

type debPkgControl struct {
	buf   *bytes.Buffer
	tw    *tar.Writer
	gw    *gzip.Writer
	info  debPkgControlInfo
	extra []string // Extra files added to the control.tar.gz. Typical usage is for conffiles, postinst, postrm, prerm.
}

// Digest file for GPG signing
type debPkgDigest struct {
	plaintext string // Plaintext package digest (empty when unsigned)
	clearsign string // GPG clearsigned package digest (empty when unsigned)
	version   int    // Always version 4 (for dpkg-sig 0.13.1+nmu2)
	signer    string // Name <email>
	date      string // Mon Jan 2 15:04:05 2006 (time.ANSIC)
	role      string // builder
	files     string // Multiple "\t<md5sum> <sha1sum> <size> <filename>"
	// E.g:
	//       3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	//       79bb73dbb522dc1a2dd1b9c2ec89fc79 26d29d15aad5c0e051d07571e28da2bc0009707e 366 control.tar.gz
	//       e1a6e48c95a760170029ef7872cec994 e02ed99e5c4fd847bde12b4c2c30dd814b26ec27 136 data.tar.gz
}

// DebPkg holds data for a single debian package
type DebPkg struct {
	debianBinary string
	control      debPkgControl
	data         debPkgData
	digest       debPkgDigest
	files        []string
}

// New creates new debian package with the following defaults:
//
//   Version: 0.0.0
func New() *DebPkg {
	d := &DebPkg{}

	d.debianBinary = debPkgDebianBinary
	d.control.info.vcsType = VcsTypeUnset
	d.control.info.priority = PriorityUnset

	d.control.buf = &bytes.Buffer{}
	d.control.gw = gzip.NewWriter(d.control.buf)
	d.control.tw = tar.NewWriter(d.control.gw)

	d.data.buf = &bytes.Buffer{}
	d.data.gw = gzip.NewWriter(d.data.buf)
	d.data.tw = tar.NewWriter(d.data.gw)

	return d
}

func (deb *DebPkg) verify() error {
	if deb.control.info.name == "" {
		return fmt.Errorf("empty package name")
	}

	if deb.control.info.architecture == "" {
		return fmt.Errorf("empty architecture")
	}

	return nil
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	err := deb.verify()
	if err != nil {
		return err
	}

	err = createControlTarGz(deb)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	// TODO move to separate function
	if filename == "" {
		filename = deb.control.info.name + "_" + deb.control.info.version.full + "_" + deb.control.info.architecture + ".deb"
	}

	return deb.createDebAr(filename)
}

// WriteSigned package with GPG entity
func (deb *DebPkg) WriteSigned(filename string, entity *openpgp.Entity, keyid string) error {
	var buf bytes.Buffer
	var cfg packet.Config
	var signer string
	cfg.DefaultHash = debPkgDigestDefaultHash

	for id := range entity.Identities {
		// TODO real search for keyid, need to investigate maybe a subkey?
		signer = id
	}

	deb.digest.date = fmt.Sprintf(time.Now().Format(time.ANSIC))
	deb.digest.signer = signer

	clearsign, err := clearsign.Encode(&buf, entity.PrivateKey, &cfg)
	if err != nil {
		return fmt.Errorf("error while signing: %s", err)
	}

	err = createControlTarGz(deb)
	if err != nil {
		return fmt.Errorf("error while creating control.tar.gz: %s", err)
	}

	deb.digest.plaintext = createDigestFileString(deb)

	if _, err = clearsign.Write([]byte(deb.digest.plaintext)); err != nil {
		return fmt.Errorf("error from Write: %s", err)
	}
	if err = clearsign.Close(); err != nil {
		return fmt.Errorf("error from Close: %s", err)
	}

	deb.digest.clearsign = buf.String()

	return deb.createDebAr(filename)
}

// SetName sets the name of the binary package (mandatory)
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Package
func (deb *DebPkg) SetName(name string) {
	deb.control.info.name = name
}

// SetVersion sets the full version string (mandatory), or user SetVersion* functions for "major.minor.patch"
// The upstream_version may contain only alphanumerics ( A-Za-z0-9 ) and the characters . + - : ~
//  (full stop, plus, hyphen, colon, tilde) and should start with a digit.
// NOTE: When the full string is set the SetVersion* function calls are ignored
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersion(version string) {
	// TODO add check for correct version string
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
	// TODO check email
	deb.control.info.maintainerEmail = email
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

// SetReplaces sets the names of packages which will be replaced. E.g: "pico"
// See:
func (deb *DebPkg) SetReplaces(replaces string) {
	deb.control.info.replaces = replaces
}

// SetHomepage sets the homepage URL of the package. E.g: "https://github.com/foo/bar"
func (deb *DebPkg) SetHomepage(url string) {
	// TODO check url
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
	deb.control.info.descr = " " + strings.Replace(descr, "\n", "\n ", -1) + "\n"
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

// AddControlExtra allows the advanced user to add custom script to the control.tar.gz Typical usage is
//  for conffiles, postinst, postrm, prerm.
// See: https://www.debian.org/doc/debian-policy/ch-maintainerscripts.html
// And: https://www.debian.org/doc/manuals/maint-guide/dother.en.html#maintscripts
func (deb *DebPkg) AddControlExtra(filename string) {
	deb.control.extra = append(deb.control.extra, filename)
}

// AddFile adds a file by filename to the package
func (deb *DebPkg) AddFile(filename string) error {
	deb.files = append(deb.files, filename)
	return nil
	/*
		fd, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fd.Close()

		stat, err := fd.Stat()
		if err != nil {
			return err
		}
		if stat.Mode().IsDir() {
			return nil
		}

		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Name = filename
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()

		// write the header to the tarball archive
		if err := deb.data.tw.WriteHeader(header); err != nil {
			return err
		}

		// copy the file data to the tarball
		if _, err := io.Copy(deb.data.tw, fd); err != nil {
			return err
		}

		// append md5sum for control.tar.gz file
		md5, _ := computeMd5(fd)
		deb.data.size += stat.Size()
		deb.data.md5sums += fmt.Sprintf("%x  %s\n", md5, filename)

		return nil
	*/
}

// AddDirectory adds a directory to the package
func (deb *DebPkg) AddDirectory(dir string) error {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if path != "." && path != ".." {
			files = append(files, path)
		}
		return nil
	})

	for _, file := range files {
		err = deb.AddFile(file)
		if err != nil {
			return err
		}
	}

	return err
}

// GetArchitecture gets the current local CPU architecture in debian-form
func GetArchitecture() string {
	arch := build.Default.GOARCH
	if arch == "386" {
		return "i386"
	}
	return arch
}

// computeMd5 from the os filedescriptor
func computeMd5(fd *os.File) (data []byte, err error) {
	var result []byte
	hash := md5.New()
	if _, err := io.Copy(hash, fd); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

// Create control file for control.tar.gz
func createControlFileString(deb *DebPkg) string {
	var o string

	// Autoset architecture
	if deb.control.info.architecture == "" {
		deb.SetArchitecture(GetArchitecture())
	}

	// Autogenerate version string (e.g "1.2.3") when unset
	if deb.control.info.version.full == "" {
		deb.control.info.version.full = fmt.Sprintf("%d.%d.%d",
			deb.control.info.version.major,
			deb.control.info.version.minor,
			deb.control.info.version.patch)
	}

	o += fmt.Sprintf("Package: %s\n", deb.control.info.name)
	o += fmt.Sprintf("Version: %s\n", deb.control.info.version.full)
	o += fmt.Sprintf("Architecture: %s\n", deb.control.info.architecture)
	o += fmt.Sprintf("Maintainer: %s <%s>\n",
		deb.control.info.maintainer,
		deb.control.info.maintainerEmail)
	o += fmt.Sprintf("Installed-Size: %d\n", deb.data.size)

	if deb.control.info.section != "" {
		o += fmt.Sprintf("Section: %s\n", deb.control.info.section)
	}
	if deb.control.info.priority != PriorityUnset {
		o += fmt.Sprintf("Priority: %s\n", deb.control.info.priority)
	}
	if deb.control.info.homepage != "" {
		o += fmt.Sprintf("Homepage: %s\n", deb.control.info.homepage)
	}
	if deb.control.info.vcsType != VcsTypeUnset && deb.control.info.vcsURL != "" {
		o += fmt.Sprintf("Vcs-%s: %s\n", deb.control.info.vcsType, deb.control.info.vcsURL)
	}
	if deb.control.info.vcsBrowser != "" {
		o += fmt.Sprintf("Vcs-Browser: %s\n", deb.control.info.vcsBrowser)
	}
	if deb.control.info.builtUsing != "" {
		o += fmt.Sprintf("Built-Using: %s\n", deb.control.info.builtUsing)
	}

	o += fmt.Sprintf("Description: %s\n", deb.control.info.descrShort)
	o += fmt.Sprintf("%s", deb.control.info.descr)

	return o
}

func digestCalcDataHash(data *bytes.Buffer, hash hash.Hash) string {
	var result []byte
	if _, err := io.Copy(hash, data); err != nil {
		return ""
	}
	return string(hash.Sum(result))
}

// Create unsigned digest file at toplevel of deb package
// NOTE: the deb.digest.version and deb.digest.role are set in this function!
func createDigestFileString(deb *DebPkg) string {
	const digestFileTmpl = `Version: %d
Signer: %s
Date: %s
Role: %s
Files: 
%s`
	deb.digest.version = debPkgDigestVersion
	deb.digest.role = debPkgDigestRole

	// debian-binary
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), md5.New()),
		digestCalcDataHash(bytes.NewBuffer([]byte(deb.debianBinary)), sha1.New()),
		len(deb.debianBinary),
		"debian-binary")

	// control.tar.gz
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		0, 0,
		len(deb.control.buf.Bytes()),
		"control.tar.gz")

	// data.tar.gz
	deb.digest.files += fmt.Sprintf("\t%x %x %d %s\n",
		0, 0,
		len(deb.data.buf.Bytes()),
		"data.tar.gz")

	return fmt.Sprintf(digestFileTmpl,
		deb.digest.version,
		deb.digest.signer,
		deb.digest.date,
		deb.digest.role,
		deb.digest.files)
}

func createControlTarGz(deb *DebPkg) error {
	body := []byte(createControlFileString(deb))
	hdr := tar.Header{
		Name:     "control",
		Size:     int64(len(body)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if err := deb.control.tw.WriteHeader(&hdr); err != nil {
		return fmt.Errorf("cannot write header of control file to control.tar.gz: %v", err)
	}

	if _, err := deb.control.tw.Write(body); err != nil {
		return fmt.Errorf("cannot write control file to control.tar.gz: %v", err)
	}

	hdr = tar.Header{
		Name:     "md5sums",
		Size:     int64(len(deb.data.md5sums)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}
	if err := deb.control.tw.WriteHeader(&hdr); err != nil {
		return fmt.Errorf("cannot write header of md5sums file to control.tar.gz: %v", err)
	}
	if _, err := deb.control.tw.Write([]byte(deb.data.md5sums)); err != nil {
		return fmt.Errorf("cannot write md5sums file to control.tar.gz: %v", err)
	}

	if err := deb.control.tw.Close(); err != nil {
		return fmt.Errorf("closing control.tar.gz: %v", err)
	}
	if err := deb.control.gw.Close(); err != nil {
		return fmt.Errorf("closing control.tar.gz: %v", err)
	}
	return nil
}

func addArFile(now time.Time, w *ar.Writer, name string, body []byte) error {
	hdr := ar.Header{
		Name:    name,
		Size:    int64(len(body)),
		Mode:    0644,
		ModTime: now,
	}
	if err := w.WriteHeader(&hdr); err != nil {
		return fmt.Errorf("cannot write file header: %v", err)
	}
	_, err := w.Write(body)
	return err
}

func (deb *DebPkg) createDebAr(filename string) error {
	// Create file
	removeDeb := true
	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create: %s", filename)
	}

	defer func() {
		fd.Close()
		if removeDeb {
			os.Remove(filename)
		}
	}()

	deb.data.tw.Close()
	deb.data.gw.Close()
	now := time.Now()
	w := ar.NewWriter(fd)

	if err := w.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("cannot write ar header to deb file: %v", err)
	}
	if err := addArFile(now, w, "debian-binary", []byte(deb.debianBinary)); err != nil {
		return fmt.Errorf("cannot pack debian-binary: %v", err)
	}
	if err := addArFile(now, w, "control.tar.gz", deb.control.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add control.tar.gz to deb: %v", err)
	}
	if err := addArFile(now, w, "data.tar.gz", deb.data.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add data.tar.gz to deb: %v", err)
	}
	if deb.digest.clearsign != "" {
		if err := addArFile(now, w, "digests.asc", []byte(deb.digest.clearsign)); err != nil {
			return fmt.Errorf("cannot add digests.asc to deb: %v", err)
		}
	}
	removeDeb = false
	return nil
}
