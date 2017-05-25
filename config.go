// Copyright 2016 Jerry Jacobs. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"fmt"
	"io/ioutil"

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
	cfg := debPkgSpecFileCfg{}

	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfgFile, &cfg)
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

	for _, file := range cfg.Files {
		err := deb.AddFile(file.Src, file.Dest)
		if err != nil {
			fmt.Printf("error adding file %s: %v\n", file.Src, err)
		}
	}

	for _, dir := range cfg.Directories {
		err := deb.AddDirectory(dir)
		if err != nil {
			fmt.Printf("error adding directory %s: %v\n", dir, err)
		}
	}

	for _, dir := range cfg.EmptyDirectories {
		err := deb.AddEmptyDirectory(dir)
		if err != nil {
			fmt.Printf("error adding directory %s: %v\n", dir, err)
		}
	}

	return nil
}
