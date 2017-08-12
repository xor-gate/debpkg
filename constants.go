// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

// Priority for Debian package
type Priority string

// Package Priority
const (
	PriorityUnset     Priority = ""          // Priority field is skipped
	PriorityRequired  Priority = "required"  // Packages which are necessary for the proper functioning of the system
	PriorityImportant Priority = "important" // Important programs, including those which one would expect to find on any Unix-like system
	PriorityStandard  Priority = "standard"  // These packages provide a reasonably small but not too limited character-mode system
	PriorityOptional  Priority = "optional"  // This is all the software that you might reasonably want to install if you didn't know what it was and don't have specialized requirements
)

// VcsType for Debian package supported version control system (Vcs) types
type VcsType string

// Package VcsType
const (
	VcsTypeUnset      VcsType = ""      // VcsType field is skipped
	VcsTypeArch       VcsType = "Arch"  // Arch
	VcsTypeBazaar     VcsType = "Bzr"   // Bazaar
	VcsTypeDarcs      VcsType = "Darcs" // Darcs
	VcsTypeGit        VcsType = "Git"   // Git
	VcsTypeMercurial  VcsType = "Hg"    // Mercurial
	VcsTypeMonotone   VcsType = "Mtn"   // Monotone
	VcsTypeSubversion VcsType = "Svn"   // Subversion
)

const (
	DefaultInstallPrefix = "/usr"
	DefaultBinDir        = "bin"
	DefaultSbinDir       = "sbin"
	DefaultSysConfDir    = "etc"
	DefaultDataRootDir   = "share"
)

const debianPathSeparator  = "/"
const debianBinaryVersion  = "2.0\n"
const debianFileExtension  = "deb"
