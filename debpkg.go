package debpkg

import (
	"archive/tar"
	"bytes"
	"time"
	"compress/gzip"
	"fmt"
)

type debPkgData struct {
	md5sums string
	buf *bytes.Buffer
	tw *tar.Writer
	gw *gzip.Writer
}

const debPkgDigestVersion = 4

// Digest file for GPG signing
type debPkgDigest struct {
	version int    // Always version 4 (for dpkg-sig 0.13.1+nmu2)
	signer  string // Name <email>
	date    string // Mon Jan 2 15:04:05 2006 (time.ANSIC)
	role    string // builder
	files   string // Multiple <md5sum> <sha1sum> <size> <filename>
                       // E.g: 
                       // 	3cf918272ffa5de195752d73f3da3e5e 7959c969e092f2a5a8604e2287807ac5b1b384ad 4 debian-binary
	               //       79bb73dbb522dc1a2dd1b9c2ec89fc79 26d29d15aad5c0e051d07571e28da2bc0009707e 366 control.tar.gz
	               //       e1a6e48c95a760170029ef7872cec994 e02ed99e5c4fd847bde12b4c2c30dd814b26ec27 136 data.tar.gz
}

type DebPkg struct {
	name            string
	version         string
	architecture    string
	installedSize   int64
	maintainer      string
	maintainerEmail string
	homepage        string
	descrShort      string
	descr           string
	data            debPkgData
	digest          debPkgDigest
}

// Create new debian package
func New() *DebPkg {
	d := &DebPkg{}

	d.data.buf = new(bytes.Buffer)
	d.data.gw  = gzip.NewWriter(d.data.buf)
	d.data.tw  = tar.NewWriter(d.data.gw)

	return d
}

// GPG sign the package
func (deb *DebPkg) Sign() {
	deb.digest.version = debPkgDigestVersion
	deb.digest.date    = fmt.Sprintf(time.Now().Format(time.ANSIC))
} 

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	fmt.Printf("%s", createControlFile(deb))
	fmt.Printf("%s", createDigestFile(deb))
	return nil
}

// Set package name
func (deb *DebPkg) SetName(name string) {
	deb.name = name
}

// Set package version
func (deb *DebPkg) SetVersion(version string) {
	deb.version = version
}

// Set maintainer. E.g: "Foo Bar"
func (deb *DebPkg) SetMaintainer(maintainer string) {
	deb.maintainer = maintainer
}

// Set maintainer email. E.g: "foo@bar.com"
func (deb *DebPkg) SetMaintainerEmail(emailAddress string) {
	// add check
	deb.maintainerEmail = emailAddress
}

// Set homepage url. E.g: "https://github.com/foo/bar"
func (deb *DebPkg) SetHomepageUrl(url string) {
	// check url
	deb.homepage = url
}

// Set short description. E.g: "My awesome foo bar baz tool"
func (deb *DebPkg) SetShortDescription(descr string) {
	deb.descrShort = descr
}

// Set long description. E.g:
// "This tool will calculation the most efficient way to world domination"
func (deb *DebPkg) SetDescription(descr string) {
	deb.descr = descr
}

func createControlFile(deb *DebPkg) string {
	const controlFileTmpl = `
Package: %s
Version: %s
Architecture: %s
Maintainer: %s <%s>
Installed-Size: %d
Section: devel
Priority: extra
Homepage: %s
Description: %s
 %s
`
	return fmt.Sprintf(controlFileTmpl,
		deb.name,
		deb.version,
		deb.architecture,
		deb.maintainer,
		deb.maintainerEmail,
		deb.installedSize,
		deb.homepage,
		deb.descrShort,
		deb.descr)
}

func createDigestFile(deb *DebPkg) string {
const digestFileTmpl = `
Version: %d
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
