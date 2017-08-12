// Copyright 2017 Debpkg authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package debpkg

import (
	"errors"
)

var ErrClosed = errors.New("debpkg: Closed")
var ErrIO     = errors.New("debpkg: I/O failed")

// setError sets the package error when not nil
// setting an error when the current error is ErrClosed it will panic
func (deb *DebPkg) setError(err error) error {
	if err == nil {
		return nil
	}
	if deb.err == ErrClosed {
		panic("debpkg: Trying to overwrite ErrClosed")
	}
	if err != ErrClosed {
		deb.err = err
	}
	return err
}
