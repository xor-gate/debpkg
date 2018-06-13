// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVarInit(t *testing.T) {
	tvs := map[string]string{
		"INSTALLPREFIX": DefaultInstallPrefix,
		"BINDIR":        DefaultBinDir,
		"SBINDIR":       DefaultSbinDir,
		"SYSCONFDIR":    DefaultSysConfDir,
		"DATAROOTDIR":   DefaultDataRootDir,
	}

	dv := DefaultVariables()

	for v, exp := range tvs {
		assert.Equal(t, exp, dv.Get(v))
	}
}

func TestGetVarWithPrefix(t *testing.T) {
	tvs := map[string]string{
		"BINDIR":      "/usr/bin",
		"SBINDIR":     "/usr/sbin",
		"SYSCONFDIR":  "/usr/etc", // FIXME should not be possible -> "/etc"
		"DATAROOTDIR": "/usr/share",
	}

	dv := DefaultVariables()
	for v, exp := range tvs {
		assert.Equal(t, exp, dv.GetWithPrefix(v))
	}
}

func TestExpandVarBinDir(t *testing.T) {
	tvs := map[string]string{
		"{{.BINDIR}}":      "/usr/bin",
		"{{.SBINDIR}}":     "/usr/sbin",
		"{{.SYSCONFDIR}}":  "/usr/etc", // FIXME should not be possible -> "/etc"
		"{{.DATAROOTDIR}}": "/usr/share",
	}

	dv := DefaultVariables()

	for val, exp := range tvs {
		res, err := dv.ExpandVar(val)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	}
}
