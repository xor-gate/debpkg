// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io/ioutil"
	"github.com/xor-gate/debpkg/internal/config"
)

// Config loads settings from a depkg.yml specfile
func (deb *DebPkg) Config(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("problem reading config file: %v", err)
	}

	dataExpanded, err := ExpandVar(string(data))
	if err != nil {
		return err
	}

	cfg, err := config.PkgSpecFileUnmarshal([]byte(dataExpanded))
	if err != nil {
		return err
	}

	deb.SetSection(cfg.Section)
	deb.SetPriority(Priority(cfg.Priority))
	deb.SetName(cfg.Name)
	deb.SetVersion(cfg.Version)
	deb.SetArchitecture(cfg.Architecture)
	deb.SetMaintainer(cfg.Maintainer)
	deb.SetMaintainerEmail(cfg.MaintainerEmail)
	deb.SetHomepage(cfg.Homepage)
	deb.SetShortDescription(cfg.Description.Short)
	deb.SetDescription(cfg.Description.Long)
	deb.SetBuiltUsing(cfg.BuiltUsing)

	for _, file := range cfg.Files {
		err := deb.AddFile(file.Src, file.Dest)
		if err != nil {
			return fmt.Errorf("error adding file %s: %v", file.Src, err)
		}
	}

	for _, dir := range cfg.Directories {
		err := deb.AddDirectory(dir)
		if err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	for _, dir := range cfg.EmptyDirectories {
		err := deb.AddEmptyDirectory(dir)
		if err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	return nil
}
