// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io/ioutil"
	"runtime"

	yaml "gopkg.in/yaml.v2"
)

type debPkgSpecFileCfg struct {
	Name            string `yaml:"name"`
	Version         string `yaml:"version"`
	Architecture    string `yaml:"architecture"`
	Maintainer      string `yaml:"maintainer"`
	MaintainerEmail string `yaml:"maintainer_email"`
	Homepage        string `yaml:"homepage"`
	Section         string `yaml:"section"`
	Priority        string `yaml:"priority"`
	BuiltUsing      string `yaml:"built_using"`
	Description     struct {
		Short string `yaml:"short"`
		Long  string `yaml:"long"`
	}
	Files []struct {
		Src  string `yaml:"file"`
		Dest string `yaml:"dest"`
	} `yaml:",flow"`
	Directories      []string `yaml:",flow"`
	EmptyDirectories []string `yaml:"emptydirs,flow"`
}

// Config loads settings from a depkg.yml specfile
func (deb *DebPkg) Config(filename string) error {
	cfg := debPkgSpecFileCfg{
		Name:            "unknown",
		Version:         "0.1.0+dev",
		Architecture:    "any",
		Maintainer:      "anonymous",
		MaintainerEmail: "anon@foo.bar",
		Homepage:        "https://www.google.com",
		Section:         "misc",
		Priority:        string(PriorityOptional),
		BuiltUsing:      runtime.Version(),
	}
	cfg.Description.Long = "-"
	cfg.Description.Short = "-"

	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("problem reading config file: %v", err)
	}
	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		return fmt.Errorf("problem unmarshaling config file: %v", err)
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
