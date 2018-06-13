// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"bytes"
	"strings"
	"text/template"
)

// Variables
type Variables map[string]string

// DefaultVariables creates a Variables object populated from Default* constants (e.g InstallPrefixVar set to DefaultInstallPrefix)
func DefaultVariables() Variables {
	v := make(Variables)
	v.Set(InstallPrefixVar, DefaultInstallPrefix)
	v.Set(BinDirVar, DefaultBinDir)
	v.Set(SbinDirVar, DefaultSbinDir)
	v.Set(SysConfDirVar, DefaultSysConfDir)
	v.Set(DataRootDirVar, DefaultDataRootDir)
	return v
}

// SetVar sets a variable for use with config file
func (v Variables) Set(key, val string) {
	v[key] = val
}

// GetVar gets a variable by key
func (v Variables) Get(key string) string {
	if val, ok := v[key]; ok {
		return val
	}
	return ""
}

// GetVarWithPrefix gets a variable and appends INSTALLPREFIX when the value doesn't start with "/"
func (v Variables) GetWithPrefix(key string) string {
	val := v.Get(key)
	if val == "" {
		return val
	}
	if strings.HasPrefix(val, debianPathSeparator) {
		return val
	}
	return v.Get(InstallPrefixVar) + debianPathSeparator + val
}

// ExpandVar expands a string with variables
func (v Variables) ExpandVar(msg string) (string, error) {
	tmpl, err := template.New("msg").Parse(msg)
	if err != nil {
		return "", err
	}
	env := struct {
		INSTALLPREFIX string
		BINDIR        string
		SBINDIR       string
		DATAROOTDIR   string
		SYSCONFDIR    string
	}{
		INSTALLPREFIX: v.Get(InstallPrefixVar),
		BINDIR:        v.GetWithPrefix(BinDirVar),
		SBINDIR:       v.GetWithPrefix(SbinDirVar),
		DATAROOTDIR:   v.GetWithPrefix(DataRootDirVar),
		SYSCONFDIR:    v.GetWithPrefix(SysConfDirVar),
	}
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, env); err != nil {
		return "", err
	}
	return buf.String(), nil
}
