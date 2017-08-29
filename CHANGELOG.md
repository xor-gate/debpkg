# Changelog

All notable changes to the [debpkg project](https://github.com/xor-gate/debpkg) will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased] - 2017-xx-xx

The first release comes with:

* Setting required fields for the control file
* Adding single files (with optional destination path)
* Adding directories with files recursivley (destination path is not supported)
* Adding empty directories
* Adding control extra files (`preinst`,`postinst`,`prerm`,`postrm`)
* Compression of the data archive with `tar.gz`
* Configuring a custom TMPDIR for intermediate files before `New` calls (defaults to OS)
* Unified configuration file and tool to generate Debian packages