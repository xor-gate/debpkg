// Copyright 2016 Jerry Jacobs. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package debpkg implements creation of (gpg signed) debian packages
//
// Overview
//
// The most basic debian package is created as follows (without error checking):
//
//  deb := debpkg.New()
//
//  deb.SetName("foobar")
//  deb.SetVersion("1.2.3")
//  deb.SetArchitecture("amd64")
//  deb.SetMaintainer("Foo Bar")
//  deb.SetMaintainerEmail("foo@bar.com")
//  deb.SetHomepage("http://foobar.com")
//
//  deb.SetShortDescription("Minimal foo bar package")
//  deb.SetDescription("Foo bar package doesn't do anything")
//
//  deb.AddFile("/tmp/foobar")
//
//  deb.Write("foobar.deb")
package debpkg
