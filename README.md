# debpkg

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Build][Build-Status-Image]][Build-Status-Url]
[![BuildAppVeyor][BuildAV-Status-Image]][BuildAV-Status-Url]
[![Coverage][Coverage-Image]][Coverage-Url]

debpkg is a pure [Go](https://golang.org) library to create [Debian](https://debian.org) Linux packages. The package can be directly imported in your go project or run a `debpkg` cli tool with a custom debpkg YAML specfile (just like a single [RPM spec file](http://ftp.rpm.org/max-rpm/s1-rpm-build-creating-spec-file.html)). It has zero dependencies outside go modules, to
 and doesn't use exec wrappers. Therefor it is able to generate packages from non Linux systems like Mac OS X, *BSD and Windows without the Perl dependency (as with official [debbuild](https://github.com/debbuild/debbuild)). The "problem" with the Debian package build system is it universal and flexable nature, which can be a pro or con. Packaging just a small bunch of files with simple configuration the Debian build system is a fairly complex beast to understand. And you must dig deep into the workings of the Debian package management documentation. Debpkg is not able to resolve all the advanced packaging requirements theirfor one must use the official `debbuild` and `debhelper`. It is created for the simple case to convert a directory structure to be installable and maintained via the standard debian way.

## Features

The feature list below describes the usability state of the project:

- Create debian packages from files and folders
- Create package from `debpkg.yml` specfile 
- Add custom control files (preinst, postinst, prerm, postrm etcetera)

It is currently not possible to use the `debpkg` as a framework to manipulate and introspect individual Debian package objects ([see issue #26](https://github.com/xor-gate/debpkg/issues/26)). As it is only capable of creating packages.

## Why this package was created

This package was created due to the lack to debianize from other platforms (windows/mac/*bsd). Because
 the whole debian package management system is glued together with Perl scripts and uses a bunch of Perl
 modules. Which isn't easy to use in a CI/CD pipeline for Go projects.

Converting a simple directory structure with files into a debian package is a difficult without knowing the `deb`-file internals.

This package is heavily inspired by [godeb](https://github.com/niemeyer/godeb) and
 [CPackDeb](https://cmake.org/cmake/help/v3.5/module/CPackDeb.html). It is very royal [licensed](LICENSE).

## Installation

`go get -u github.com/xor-gate/debpkg/cmd/debpkg`

## Status

The package is currently in production state. The API is still unstable (v0) most rough edges are already polished 
 but a wider audience is necessary. 

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

This project originate from an in-company implementation sponsored by [@dualinventive](https://github.com/dualinventive) in 2016-2017, with help from collegue [@rikvdh](https://github.com/rikvdh).

I would like to mention some other noticable projects for debian package management:

* https://github.com/Debian/dh-make-golang
* https://github.com/niemeyer/godeb
* https://github.com/smira/aptly
* https://github.com/esell/deb-simple
* https://github.com/paultag/go-debian
* https://github.com/jordansissel/fpm
* https://github.com/laher/debgo-v0.2
* https://github.com/debber/debber-v0.3
* https://github.com/laher/goxc
* https://github.com/knqyf263/go-deb-version

Other packaging tools in golang:

* https://github.com/google/rpmpack

## Debugging

When debugging on a Debian system the following commands are usefull:

* Print package info: `dpkg --info <debfile>`
* Extract data.tar.gz: `dpkg --extract <debfile> data`
* Verbose extract data.tar.gz: `dpkg --vextract <debfile> data`
* Extract control.tar.gz: `dpkg --control <debfile> control`

## Further documentation

More information can be gathered from the Debian and Ubuntu documentation:

* [dpkg manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg)
* [dpkg-deb manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg)
* [dpkg-sig manpage](https://manpages.debian.org/cgi-bin/man.cgi?query=dpkg-sig)
* [Debian New Maintainers' Guide](https://www.debian.org/doc/manuals/maint-guide/)
* [Debian Policy Manual](https://www.debian.org/doc/debian-policy/)
* [Setting up signed APT repository with Reprepro](https://wiki.debian.org/SettingUpSignedAptRepositoryWithReprepro)
* [Create authenticated repository](https://help.ubuntu.com/community/CreateAuthenticatedRepository)

## License

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
