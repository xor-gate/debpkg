// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestVarInit(t *testing.T) {
	tvs := map[string]string {
		"INSTALLPREFIX": DefaultInstallPrefix,
		"BINDIR":        DefaultBinDir,
		"SBINDIR":       DefaultSbinDir,
		"SYSCONFDIR":    DefaultSysConfDir,
		"DATAROOTDIR":   DefaultDataRootDir,
	}

	for v, exp := range tvs {
		assert.Equal(t, exp, GetVar(v))
	}
}

func TestGetVarWithPrefix(t *testing.T) {
	tvs := map[string]string {
		"BINDIR":      "/usr/bin",
		"SBINDIR":     "/usr/sbin",
		"SYSCONFDIR":  "/usr/etc", // FIXME should not be possible -> "/etc"
		"DATAROOTDIR": "/usr/share",
	}

	for v, exp := range tvs {
		assert.Equal(t, exp, GetVarWithPrefix(v))
	}
}

func TestExpandVarBinDir(t *testing.T) {
	tvs := map[string]string {
		"{{.BINDIR}}":      "/usr/bin",
		"{{.SBINDIR}}":     "/usr/sbin",
		"{{.SYSCONFDIR}}":  "/usr/etc", // FIXME should not be possible -> "/etc"
		"{{.DATAROOTDIR}}": "/usr/share",
	}

	for val, exp := range tvs {
		res, err := ExpandVar(val)
		assert.Nil(t, err)
		assert.Equal(t, exp, res)
	}
}
