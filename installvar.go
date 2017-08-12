// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"bytes"
	"strings"
	"text/template"
)

var vars map[string]string

func init() {
	vars = make(map[string]string)
	vars["INSTALLPREFIX"] = DefaultInstallPrefix
	vars["BINDIR"]        = DefaultBinDir
	vars["SBINDIR"]       = DefaultSbinDir
	vars["SYSCONFDIR"]    = DefaultSysConfDir
	vars["DATAROOTDIR"]   = DefaultDataRootDir
}

// SetVar sets a variable for use with config file
func SetVar(key, val string) {
	vars[key] = val
}

// GetVar gets a variable
func GetVar(v string) string {
	if val, ok := vars[v]; ok {
		return val
	}
	return ""
}

// GetVarWithPrefix gets a variable and appends INSTALLPREFIX when the value doesn't start with "/"
func GetVarWithPrefix(v string) string{
	val := GetVar(v)
	if val == "" {
		return val
	}
	if strings.HasPrefix(val, debianPathSeparator) {
		return val
	}
	return vars["INSTALLPREFIX"] + debianPathSeparator + val
}

// ExpandVar expands a string with variables
func ExpandVar(msg string) (string, error) {
	tmpl, err := template.New("msg").Parse(msg)
	if err != nil {
		return "", err
	}
	env := struct {
		INSTALLPREFIX string
		BINDIR string
		SBINDIR string
		DATAROOTDIR string
		SYSCONFDIR string
	}{
		INSTALLPREFIX: vars["INSTALLPREFIX"],
		BINDIR:        GetVarWithPrefix("BINDIR"),
		SBINDIR:       GetVarWithPrefix("SBINDIR"),
		DATAROOTDIR:   GetVarWithPrefix("DATAROOTDIR"),
		SYSCONFDIR:    GetVarWithPrefix("SYSCONFDIR"),
	}
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, env); err != nil {
		return "",err
	}
	return buf.String(),nil
}
