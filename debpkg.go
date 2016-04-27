package debpkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/blakesmith/ar"
	"github.com/go-yaml/yaml"
	"go/build"
	"io"
	"strings"
	"os"
	"path/filepath"
	"time"
)

type VcsType string

const (
	VcsTypeArch       VcsType = "Arch"
	VcsTypeBazaar     VcsType = "Bzr"
	VcsTypeDarcs      VcsType = "Darcs"
	VcsTypeGit        VcsType = "Git"
	VcsTypeMercurial  VcsType = "Hg"
	VcsTypeMonotone   VcsType = "Mtn"
	VcsTypeSubversion VcsType = "Svn"
)

const debPkgDebianBinaryVersion = "2.0\n"
const debPkgDigestVersion = 4
const debPkgDigestRole = "builder"

type debPkgSpecFileCfg struct {
	Description struct {
		Short string `yaml:"short"`
		Long  string `yaml:"long"`
	}
}

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
	priority        string
	descrShort      string  // Short package description
	descr           string  // Long package description
	vcsType         VcsType // E.g: "Svn", "Git" etcetera
	vcsURL          string  // E.g: git@github.com:xor-gate/debpkg.git
	vcsBrowser      string  // E.g: https://github.com/xor-gate/debpkg
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
	version int    // Always version 4 (for dpkg-sig 0.13.1+nmu2)
	signer  string // Name <email>
	date    string // Mon Jan 2 15:04:05 2006 (time.ANSIC)
	role    string // builder
	files   string // Multiple <md5sum> <sha1sum> <size> <filename>
	// E.g:
	//       3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	//       79bb73dbb522dc1a2dd1b9c2ec89fc79 26d29d15aad5c0e051d07571e28da2bc0009707e 366 control.tar.gz
	//       e1a6e48c95a760170029ef7872cec994 e02ed99e5c4fd847bde12b4c2c30dd814b26ec27 136 data.tar.gz
}

type DebPkg struct {
	control debPkgControl
	data    debPkgData
	digest  debPkgDigest
}

// Create new debian package
func New() *DebPkg {
	d := &DebPkg{}

	d.control.buf = &bytes.Buffer{}
	d.control.gw = gzip.NewWriter(d.control.buf)
	d.control.tw = tar.NewWriter(d.control.gw)

	d.data.buf = &bytes.Buffer{}
	d.data.gw = gzip.NewWriter(d.data.buf)
	d.data.tw = tar.NewWriter(d.data.gw)

	return d
}

// GPG sign the package
func (deb *DebPkg) Sign() {
	deb.digest.version = debPkgDigestVersion
	deb.digest.date = fmt.Sprintf(time.Now().Format(time.ANSIC))
	deb.digest.role = debPkgDigestRole
}

// Load configuration from depkg.yml specfile
func (deb *DebPkg) Config(filename string) error {
	cfg := debPkgSpecFileCfg{}
	data := new(bytes.Buffer)

	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	io.Copy(data, fd)

	err = yaml.Unmarshal(data.Bytes(), &cfg)
	if err != nil {
		return err
	}

	fmt.Printf("config:\n\n%+v\n", cfg)
	return nil
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer fd.Close()

	deb.data.tw.Close()
	deb.data.gw.Close()

	createControlTarGz(deb)
	deb.createDebAr(fd)

	return nil
}

// Set package name (mandatory)
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Package
func (deb *DebPkg) SetName(name string) {
	deb.control.info.name = name
}

// Set package version string (mandatory)
// NOTE: When the full string is set the SetVersion* function calls are ignored
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersion(version string) {
	deb.control.info.version.full = version
}

// Set package version major number
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionMajor(major uint) {
	deb.control.info.version.major = major
}

// Set package version minor number
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionMinor(minor uint) {
	deb.control.info.version.minor = minor
}

// Set package version patch number
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version
func (deb *DebPkg) SetVersionPatch(patch uint) {
	deb.control.info.version.patch = patch
}

// Set architecture. E.g "i386, amd64, arm". See `dpkg-architecture -L` for all supported.
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

// Set maintainer (mandatory). E.g: "Foo Bar"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (deb *DebPkg) SetMaintainer(maintainer string) {
	deb.control.info.maintainer = maintainer
}

// Set maintainer email. E.g: "foo@bar.com"
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Maintainer
func (deb *DebPkg) SetMaintainerEmail(email string) {
	// add check
	deb.control.info.maintainerEmail = email
}

// Set suggests. E.g: aptitude
func (deb *DebPkg) SetSuggests(suggests string) {
	deb.control.info.suggests = suggests
}

// Set conflicts. E.g: nano
func (deb *DebPkg) SetConflicts(conflicts string) {
	deb.control.info.conflicts = conflicts
}

// Set provides. E.g: editor
func (deb *DebPkg) SetProvides(provides string) {
	deb.control.info.provides = provides
}

// Set priority (recommended). E.g: important
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Priority
func (deb *DebPkg) SetPriority(prio string) {
	deb.control.info.priority = prio
}

// Set section (recommended). E.g: editors
// See: https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Section
func (deb *DebPkg) SetSection(section string) {
	deb.control.info.section = section
}

// Set replaces. E.g: pico
func (deb *DebPkg) SetReplaces(replaces string) {
	deb.control.info.replaces = replaces
}

// Set homepage url. E.g: "https://github.com/foo/bar"
func (deb *DebPkg) SetHomepage(url string) {
	// check url
	deb.control.info.homepage = url
}

// Set short description. E.g: "My awesome foo bar baz tool"
func (deb *DebPkg) SetShortDescription(descr string) {
	deb.control.info.descrShort = descr
}

// Set long description. E.g:
// "This tool will calculation the most efficient way to world domination"
// NOTE: The debian control file has a special formatting of the long description
//        this function replaces newlines with a newline and a space.
func (deb *DebPkg) SetDescription(descr string) {
	deb.control.info.descr = " " + strings.Replace(descr, "\n", "\n ", -1) + "\n"
}

// Set version control system (Vcs) type.
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsType(vcs VcsType) {
	deb.control.info.vcsType = vcs
}

// Set version control system (Vcs) URL
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsURL(url string) {
	deb.control.info.vcsURL = url
}

// Set version control system (Vcs) browsable source-tree URL
// See: https://www.debian.org/doc/manuals/developers-reference/best-pkging-practices.html#s6.2.5.2
func (deb *DebPkg) SetVcsBrowser(url string) {
	deb.control.info.vcsBrowser = url
}

// Allow advanced user to add custom script to the control.tar.gz Typical usage is for
//  conffiles, postinst, postrm, prerm.
func (deb *DebPkg) AddControlExtra(filename string) {
	deb.control.extra = append(deb.control.extra, filename)
}

func (deb *DebPkg) AddFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
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
		if _, err := io.Copy(deb.data.tw, file); err != nil {
			return err
		}

		md5, size, _ := computeMd5(filename)
		deb.data.size += size
		deb.data.md5sums += fmt.Sprintf("%x  %s\n", md5, filename)
	}
	return nil
}

func (deb *DebPkg) AddDirectory(dir string) error {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	for _, file := range files {
		deb.AddFile(file)
	}

	return err
}

// Get current architecture of build
func GetArchitecture() string {
	arch := build.Default.GOARCH
	if arch == "386" {
		return "i386"
	}
	return arch
}

func computeMd5(filePath string) (data []byte, size int64, err error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, 0, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, 0, err
	}

	fi, err := file.Stat()
	if err != nil {
		return result, 0, err
	}

	return hash.Sum(result), fi.Size(), nil
}

// Create control file for control.tar.gz
func createControlFileString(deb *DebPkg) string {
	var o string

	// TODO we should check if all recommended fields exist with a verify function

	if deb.control.info.architecture == "" {
		deb.SetArchitecture(GetArchitecture())
	}

	if deb.control.info.version.full == "" {
		deb.control.info.version.full = fmt.Sprintf("%d.%d.%d", deb.control.info.version.major, deb.control.info.version.minor, deb.control.info.version.patch)
	}

	o += fmt.Sprintf("Package: %s\n", deb.control.info.name)
	o += fmt.Sprintf("Version: %s\n", deb.control.info.version.full)
	o += fmt.Sprintf("Architecture: %s\n", deb.control.info.architecture)
	o += fmt.Sprintf("Maintainer: %s <%s>\n", deb.control.info.maintainer, deb.control.info.maintainerEmail)
	o += fmt.Sprintf("Installed-Size: %d\n", deb.data.size)
	// TODO Section suggested?, check with docs
	o += fmt.Sprintf("Section: %s\n", deb.control.info.section)
	// TODO Priority suggested? check with docs
	o += fmt.Sprintf("Priority: %s\n", deb.control.info.priority)
	o += fmt.Sprintf("Homepage: %s\n", deb.control.info.homepage)
	o += fmt.Sprintf("Description: %s\n", deb.control.info.descrShort)
	o += fmt.Sprintf("%s", deb.control.info.descr)

	return o
}

// Create unsigned digest file at toplevel of deb package
func createDigestFile(deb *DebPkg) string {
	const digestFileTmpl = `Version: %d
Date: %s
Signer: %s
Role: %s 
`
	return fmt.Sprintf(digestFileTmpl,
		deb.digest.version,
		deb.digest.date,
		deb.digest.signer,
		deb.digest.role)
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

func (deb *DebPkg) createDebAr(dst io.Writer) error {
	now := time.Now()
	w := ar.NewWriter(dst)
	if err := w.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("cannot write ar header to deb file: %v", err)
	}
	if err := addArFile(now, w, "debian-binary", []byte(debPkgDebianBinaryVersion)); err != nil {
		return fmt.Errorf("cannot pack debian-binary: %v", err)
	}
	if err := addArFile(now, w, "control.tar.gz", deb.control.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add control.tar.gz to deb: %v", err)
	}
	if err := addArFile(now, w, "data.tar.gz", deb.data.buf.Bytes()); err != nil {
		return fmt.Errorf("cannot add data.tar.gz to deb: %v", err)
	}
	return nil
}
