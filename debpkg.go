package debpkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
)

type debPkgData struct {
	md5sums string
	buf *bytes.Buffer
	tw *tar.Writer
	gw *gzip.Writer
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
}

// Create new debian package
func New() *DebPkg {
	d := &DebPkg{}
	d.data.buf = new(bytes.Buffer)
	d.data.gw  = gzip.NewWriter(d.data.buf)
	d.data.tw  = tar.NewWriter(d.data.gw)
	return d
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	fmt.Printf("%s", createControlFile(deb))
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
