// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io/ioutil"
	"strings"

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
		if len(file.Src) > 0 {
			if err := deb.AddFile(file.Src, file.Dest); err != nil {
				return fmt.Errorf("error adding file %s: %v", file.Src, err)
			}
		} else if len(file.Content) > 0 {
			if err := deb.AddFileString(file.Content, file.Dest); err != nil {
				return fmt.Errorf("error adding file by string: %v", err)
			}
		} else {
			return fmt.Errorf("need either 'content' or a 'src' to add a file")
		}
	}

	for _, dir := range cfg.Directories {
		if err := deb.AddDirectory(dir); err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	for _, dir := range cfg.EmptyDirectories {
		err := deb.AddEmptyDirectory(dir)
		if err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	if len(cfg.ControlExtra.Preinst) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Preinst, "\n") {
			deb.AddControlExtraString("preinst", cfg.ControlExtra.Preinst)
		} else {
			deb.AddControlExtra("preinst", cfg.ControlExtra.Preinst)
		}
	}

	if len(cfg.ControlExtra.Postinst) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Postinst, "\n") {
			deb.AddControlExtraString("postinst", cfg.ControlExtra.Postinst)
		} else {
			deb.AddControlExtra("postinst", cfg.ControlExtra.Postinst)
		}
	}

	if len(cfg.ControlExtra.Prerm) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Prerm, "\n") {
			deb.AddControlExtraString("prerm", cfg.ControlExtra.Prerm)
		} else {
			deb.AddControlExtra("prerm", cfg.ControlExtra.Prerm)
		}
	}

	if len(cfg.ControlExtra.Postrm) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Postrm, "\n") {
			deb.AddControlExtraString("postrm", cfg.ControlExtra.Postrm)
		} else {
			deb.AddControlExtra("postrm", cfg.ControlExtra.Postrm)
		}
	}
	return nil
}
