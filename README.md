# debpkg

Golang package for creating (signed) debian packages

[![License][License-Image]][License-Url] ![Stability][Stability-Status-Image] [![Godoc][Godoc-Image]][Godoc-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build][Build-Status-Image]][Build-Status-Url] [![Coverage][Coverage-Image]][Coverage-Url]

**Features**

- [ ] Create simple debian packages from files and folders
- [ ] Add custom control files (preinst, postinst, prerm, postrm etcetera)
- [ ] [dpkg](http://manpages.ubuntu.com/manpages/precise/man1/dpkg.1.html) like tool with a subset of commands (--contents, --control, --extract, --info)
- [ ] Create package from debpkg.yml specfile (like [packager.io](https://packager.io) without cruft)
- [ ] GPG sign package
- [ ] GPG verify package

# debpkg.yml specfile

The specfile is written in the [YAML markup format](http://yaml.org/). It controls
 the package information and data.

A simple example is listed below:

```
# debpkg.yml specfile
description:
  - version: 0.0.1
  - short: This package is just a test
  - long:
      This package tests the working of debpkg.
      And can wrap multiple
      lines.

      And multiple paragraphs.
```

# License

[MIT](LICENSE)

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[Stability-Status-Image]: http://badges.github.io/stability-badges/dist/experimental.svg
[Build-Status-Url]: http://travis-ci.org/xor-gate/debpkg
[Build-Status-Image]: https://travis-ci.org/xor-gate/debpkg.svg?branch=master
[Godoc-Url]: https://godoc.org/github.com/xor-gate/debpkg
[Godoc-Image]: https://godoc.org/github.com/xor-gate/debpkg?status.svg
[ReportCard-Url]: http://goreportcard.com/report/xor-gate/debpkg
[ReportCard-Image]: http://goreportcard.com/badge/xor-gate/debpkg
[Coverage-Url]: https://coveralls.io/r/xor-gate/debpkg?branch=master
[Coverage-image]: https://img.shields.io/coveralls/xor-gate/debpkg.svg
