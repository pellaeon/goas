// Tideland Go Application Support - Logger
//
// Copyright (C) 2012-2015 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Logger of the Tideland Go Application Support provides a flexible
// way to log information with different levels and on different backends.
package logger

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/pellaeon/goas/v1/version"
)

//--------------------
// VERSION
//--------------------

// PackageVersion returns the version of the version package.
func PackageVersion() version.Version {
	return version.New(3, 0, 2)
}

// EOF
