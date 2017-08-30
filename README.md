# debpkg

debpkg is a pure [Go](https://golang.org) library to create Debian packages. It has zero dependencies to
 Debian. It is able to generate packages from Mac OS X, *BSD and Windows. 

[![License][License-Image]][License-Url]
![Stability][Stability-Status-Image]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Build][Build-Status-Image]][Build-Status-Url]
[![BuildAppVeyor][BuildAV-Status-Image]][BuildAV-Status-Url]
[![Coverage][Coverage-Image]][Coverage-Url]

**Features**

The feature list below describes the usability state of the project:

- [x] Create simple debian packages from files and folders
- [x] Create package from debpkg.yml specfile 
- [x] Add custom control files (preinst, postinst, prerm, postrm etcetera)
- [ ] [dpkg](http://manpages.ubuntu.com/manpages/precise/man1/dpkg.1.html) like tool with a subset of commands (--info, --extract, --vextract, --control)
- [ ] GPG sign package (implemented but doesnt work yet with `dpkg -i`)

## Why this package was created

This package was created due to the lack to debianize from other platforms (windows/mac/*bsd). Because
 the whole debian package management system is glued together with Perl scripts and uses a bunch of Perl
 modules.

And converting a directory and files into a debian package is a pain without knowing the `deb`-file internals.

This package is heavily inspired by [godeb](https://github.com/niemeyer/godeb) and
 [CPackDeb](https://cmake.org/cmake/help/v3.5/module/CPackDeb.html). It is very royal [licensed](LICENSE).

## Installation

`go get -u github.com/xor-gate/debpkg/cmd/debpkg`

## Status

The package is currently in experimental state. The API is unstable and it has not throughout
 been tested and many things are unfinished or not implemented at all. **USE AT YOUR OWN RISK!**

# debpkg.yml specfile

The specfile is written in the [YAML markup format](http://yaml.org/). It controls
 the package information and data.

A simple example is listed below:

```
# debpkg.yml specfile
name: foobar
version: 7.6.5
architecture: all
maintainer: Foo Bar
maintainer_email: foo@bar.com
homepage: https://github.com/xor-gate/debpkg
description:
  short: This package is just a test
  long: >
    This package tests the working of debpkg.
    And can wrap multiple
    lines.

    And multiple paragraphs.
```

# Mentions

I would like to mention some other noticable projects:

* https://github.com/Debian/dh-make-golang
* https://github.com/niemeyer/godeb
* https://github.com/smira/aptly
* https://github.com/esell/deb-simple
* https://github.com/paultag/go-debian
* https://github.com/jordansissel/fpm
* https://github.com/laher/debgo-v0.2
* https://github.com/debber/debber-v0.3
* https://github.com/laher/goxc

# Debugging

* Print package info: `dpkg --info <debfile>`
* Extract data.tar.gz: `dpkg --extract <debfile> data`
* Verbose extract data.tar.gz: `dpkg --vextract <debfile> data`
* Extract control.tar.gz: `dpkg --control <debfile> control`

# Further documentation

* [dpkg manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg)
* [dpkg-deb manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg)
* [dpkg-sig manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg-sig)
* [Debian New Maintainers' Guide](https://www.debian.org/doc/manuals/maint-guide/)
* [Debian Policy Manual](https://www.debian.org/doc/debian-policy/)
* [Setting up signed APT repository with Reprepro](https://wiki.debian.org/SettingUpSignedAptRepositoryWithReprepro)
* [Create authenticated repository](https://help.ubuntu.com/community/CreateAuthenticatedRepository)

# License

[MIT](LICENSE)

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[Stability-Status-Image]: http://badges.github.io/stability-badges/dist/experimental.svg
[Build-Status-Url]: http://travis-ci.org/xor-gate/debpkg
[Build-Status-Image]: https://travis-ci.org/xor-gate/debpkg.svg?branch=master
[BuildAV-Status-Url]: https://ci.appveyor.com/project/xor-gate/debpkg
[BuildAV-Status-Image]: https://ci.appveyor.com/api/projects/status/iuw1j84l33ynxs32?svg=true
[Godoc-Url]: https://godoc.org/github.com/xor-gate/debpkg
[Godoc-Image]: https://godoc.org/github.com/xor-gate/debpkg?status.svg
[ReportCard-Url]: http://goreportcard.com/report/xor-gate/debpkg
[ReportCard-Image]: https://goreportcard.com/badge/github.com/xor-gate/debpkg
[Coverage-Url]: https://coveralls.io/r/xor-gate/debpkg?branch=master
[Coverage-image]: https://img.shields.io/coveralls/xor-gate/debpkg.svg
