package debpkg

import (
	"archive/tar"
	"bytes"
	"time"
	"go/build"
	"crypto/md5"
	"compress/gzip"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"github.com/blakesmith/ar"
)

const debPkgDebianBinaryVersion = "2.0\n"
const debPkgDigestVersion = 4
const debPkgDigestRole    = "builder"

type debPkgData struct {
	size int64
	md5sums string
	buf *bytes.Buffer
	tw *tar.Writer
	gw *gzip.Writer
}

type debPkgControlInfo struct {
	name            string
	version         string
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
	descrShort      string
	descr           string
}

type debPkgControl struct {
	buf *bytes.Buffer
	tw *tar.Writer
	gw *gzip.Writer
	info debPkgControlInfo
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
	control  debPkgControl
	data     debPkgData
	digest   debPkgDigest
}

// Create new debian package
func New() *DebPkg {
	d := &DebPkg{}

	d.control.buf = &bytes.Buffer{}
	d.control.gw  = gzip.NewWriter(d.control.buf)
	d.control.tw  = tar.NewWriter(d.control.gw)

	d.data.buf = &bytes.Buffer{}
	d.data.gw  = gzip.NewWriter(d.data.buf)
	d.data.tw  = tar.NewWriter(d.data.gw)

	return d
}

// GPG sign the package
func (deb *DebPkg) Sign() {
	deb.digest.version = debPkgDigestVersion
	deb.digest.date    = fmt.Sprintf(time.Now().Format(time.ANSIC))
	deb.digest.role    = debPkgDigestRole
}

// Write the debian package to the filename
func (deb *DebPkg) Write(filename string) error {
	fmt.Printf("control:\n\n%s\n", createControlFile(deb))
	fmt.Printf("control md5sums:\n\n%s\n", deb.data.md5sums)
	fmt.Printf("digest:\n\n%s\n", createDigestFile(deb))

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

// Set package name
func (deb *DebPkg) SetName(name string) {
	deb.control.info.name = name
}

// Set package version
func (deb *DebPkg) SetVersion(version string) {
	deb.control.info.version = version
}

// Set architecture
func (deb *DebPkg) SetArchitecture(arch string) {
	deb.control.info.architecture = arch
}

// Set maintainer. E.g: "Foo Bar"
func (deb *DebPkg) SetMaintainer(maintainer string) {
	deb.control.info.maintainer = maintainer
}

// Set maintainer email. E.g: "foo@bar.com"
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

// Set priority. E.g: important
func (deb *DebPkg) SetPriority(prio string) {
	deb.control.info.priority = prio
}

// Set section. E.g: editors
func (deb *DebPkg) SetSection(section string) {
	deb.control.info.section = section
}

// Set replaces. E.g: pico
func (deb *DebPkg) SetReplaces(replaces string) {
	deb.control.info.replaces = replaces
}

// Set homepage url. E.g: "https://github.com/foo/bar"
func (deb *DebPkg) SetHomepageUrl(url string) {
	// check url
	deb.control.info.homepage = url
}

// Set short description. E.g: "My awesome foo bar baz tool"
func (deb *DebPkg) SetShortDescription(descr string) {
	deb.control.info.descrShort = descr
}

// Set long description. E.g:
// "This tool will calculation the most efficient way to world domination"
func (deb *DebPkg) SetDescription(descr string) {
	deb.control.info.descr = descr
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
		fmt.Println(file)
		deb.AddFile(file)
	}

	return err
}

// Get debianized current architecture
func GetCurrentArchitecture() string {
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
func createControlFile(deb *DebPkg) string {
	const controlFileTmpl = `Package: %s
Version: %s
Architecture: %s
Maintainer: %s <%s>
Installed-Size: %d
Section: %s
Priority: %s
Homepage: %s
Description: %s
 %s
`
	if (deb.control.info.architecture == "") {
		deb.SetArchitecture(GetCurrentArchitecture())
	}

	return fmt.Sprintf(controlFileTmpl,
		deb.control.info.name,
		deb.control.info.version,
		deb.control.info.architecture,
		deb.control.info.maintainer,
		deb.control.info.maintainerEmail,
		deb.data.size,
		deb.control.info.section,
		deb.control.info.priority,
		deb.control.info.homepage,
		deb.control.info.descrShort,
		deb.control.info.descr)
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
	body := []byte(createControlFile(deb))
	hdr := tar.Header{
		Name:     "control",
		Size:     int64(len(body)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}
	fmt.Println(hdr)
	fmt.Printf("tw: %p", deb.control.tw)
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
