package debpkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xor-gate/debpkg/internal/targzip"
)

func newData(t *testing.T) *data {
	tgz, err := targzip.NewTempFile(os.TempDir())
	assert.Nil(t, err)
	return &data{tgz: tgz}
}

func TestDataAddDirectory(t *testing.T) {
	d := newData(t)

	err := d.addDirectory("/my/foo/directory/")
	assert.Nil(t, err)
	assert.Equal(t, []string{"/my", "/my/foo", "/my/foo/directory"}, d.dirs)
	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
}

func TestDataAddDirectoryError(t *testing.T) {
	d := newData(t)
	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
	err := d.addDirectory("/doesnt/matter")
	assert.NotNil(t, err)
}

func TestDataAddDirectoryCwd(t *testing.T) {
	d := newData(t)
	err := d.addDirectory(".")
	assert.Nil(t, err)
	assert.Empty(t, d.dirs)
	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
}

func TestDataAddFileString(t *testing.T) {
	d := newData(t)
	err := d.addFileString("test", "/foo")
	assert.Nil(t, err)
	assert.Empty(t, d.dirs)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6  /foo\n", d.md5sums)

	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
}

func TestDataAddFileStringError(t *testing.T) {
	d := newData(t)
	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
	err := d.addFileString("test", "/foo")
	assert.NotNil(t, err)
	assert.Empty(t, d.md5sums)
}

func TestDataAddFileError(t *testing.T) {
	d := newData(t)
	assert.Nil(t, d.tgz.Close())
	assert.Nil(t, os.Remove(d.tgz.Name()))
	err := d.addFileString("data_test.go", "/foo/bar.go")
	assert.NotNil(t, err)
	assert.Empty(t, d.md5sums)
}
