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
func (pkg *Package) Config(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("problem reading config file: %v", err)
	}

	dataExpanded, err := pkg.Variables.ExpandVar(string(data))
	if err != nil {
		return err
	}

	cfg, err := config.PkgSpecFileUnmarshal([]byte(dataExpanded))
	if err != nil {
		return err
	}

	pkg.SetSection(cfg.Section)
	pkg.SetPriority(Priority(cfg.Priority))
	pkg.SetName(cfg.Name)
	pkg.SetVersion(cfg.Version)
	pkg.SetArchitecture(cfg.Architecture)
	pkg.SetMaintainer(cfg.Maintainer)
	pkg.SetMaintainerEmail(cfg.MaintainerEmail)
	pkg.SetHomepage(cfg.Homepage)
	pkg.SetShortDescription(cfg.Description.Short)
	pkg.SetDescription(cfg.Description.Long)
	pkg.SetBuiltUsing(cfg.BuiltUsing)
	pkg.SetDepends(cfg.Depends)
	pkg.SetRecommends(cfg.Recommends)
	pkg.SetSuggests(cfg.Suggests)
	pkg.SetConflicts(cfg.Conflicts)
	pkg.SetProvides(cfg.Provides)
	pkg.SetReplaces(cfg.Replaces)

	for _, file := range cfg.Files {
		if len(file.File) > 0 {
			if err := pkg.AddFile(file.File, file.Dest); err != nil {
				return fmt.Errorf("error adding file %s: %v", file.File, err)
			}
		} else if len(file.Content) > 0 {
			if err := pkg.AddFileString(file.Content, file.Dest); err != nil {
				return fmt.Errorf("error adding file by string: %v", err)
			}
		} else {
			return fmt.Errorf("need either 'content' or a 'src' to add a file")
		}
		if file.ConfigFile {
			pkg.MarkConfigFile(file.Dest)
		}
	}

	for _, dir := range cfg.Directories {
		if err := pkg.AddDirectory(dir); err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	for _, dir := range cfg.EmptyDirectories {
		err := pkg.AddEmptyDirectory(dir)
		if err != nil {
			return fmt.Errorf("error adding directory %s: %v", dir, err)
		}
	}

	if len(cfg.ControlExtra.Preinst) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Preinst, "\n") {
			pkg.AddControlExtraString("preinst", cfg.ControlExtra.Preinst)
		} else {
			pkg.AddControlExtra("preinst", cfg.ControlExtra.Preinst)
		}
	}

	if len(cfg.ControlExtra.Postinst) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Postinst, "\n") {
			pkg.AddControlExtraString("postinst", cfg.ControlExtra.Postinst)
		} else {
			pkg.AddControlExtra("postinst", cfg.ControlExtra.Postinst)
		}
	}

	if len(cfg.ControlExtra.Prerm) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Prerm, "\n") {
			pkg.AddControlExtraString("prerm", cfg.ControlExtra.Prerm)
		} else {
			pkg.AddControlExtra("prerm", cfg.ControlExtra.Prerm)
		}
	}

	if len(cfg.ControlExtra.Postrm) > 0 {
		if strings.ContainsAny(cfg.ControlExtra.Postrm, "\n") {
			pkg.AddControlExtraString("postrm", cfg.ControlExtra.Postrm)
		} else {
			pkg.AddControlExtra("postrm", cfg.ControlExtra.Postrm)
		}
	}
	return nil
}
