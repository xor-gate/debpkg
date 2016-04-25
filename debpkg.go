package debpkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
)

type debpkgData struct {
	md5sums string
	buf *bytes.Buffer
	tw *tar.Writer
	gw *gzip.Writer
}

type debpkg struct {
	name            string
	version         string
	architecture    string
	installedSize   int64
	maintainer      string
	maintainerEmail string
	homepage        string
	descrShort      string
	descr           string
	data debpkgData
}

func New() *debpkg {
	d := &debpkg{}
	d.data.buf = new(bytes.Buffer)
	d.data.gw  = gzip.NewWriter(d.data.buf)
	d.data.tw  = tar.NewWriter(d.data.gw)
	return d
}

func (deb *debpkg) SetName(name string) {
	deb.name = name
}

func (deb *debpkg) SetVersion(version string) {
	deb.version = version
}

func (deb *debpkg) SetMaintainer(maintainer string) {
	deb.maintainer = maintainer
}

func (deb *debpkg) SetMaintainerEmail(emailAddress string) {
	// add check
	deb.maintainerEmail = emailAddress
}

func (deb *debpkg) SetHomepageUrl(url string) {
	// check url
	deb.homepage = url
}

func (deb *debpkg) SetShortDescription(descr string) {
	deb.descrShort = descr
}

func (deb *debpkg) SetDescription(descr string) {
	deb.descr = descr
}

func createControlFile(deb *debpkg) string {
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

func (deb *debpkg) Write(filename string) error {
	fmt.Printf("%s", createControlFile(deb))
	return nil
}
