// Copyright (c) 2021 - for information on the respective copyright owner
// see the NOTICE file and/or the repository <FIXME-repository-address>.
//
// SPDX-License-Identifier: Apache-2.0

package assets2036go

import "log"

type logger struct {
	_log bool
}

// one instance
var _logger *logger

func logg() *logger {
	if _logger == nil {
		_logger = &logger{
			_log: true,
		}
	}
	return _logger
}

// SwitchOnLogging switches on or off logging for the assets2036 lib
func SwitchOnLogging(on bool) {
	_logger._log = on
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func (logger *logger) Printf(format string, v ...interface{}) {
	if logger._log {
		log.Printf(format, v...)
	}
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func (logger *logger) Print(err error) {
	if logger._log {
		log.Print(err)
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func (logger *logger) Println(v ...interface{}) {
	if logger._log {
		log.Println(v...)
	}
}

func (logger *logger) Fatal(v ...interface{}) {
	if logger._log {
		log.Fatal(v...)
	}
}

func (logger *logger) Fatalf(format string, v ...interface{}) {
	if logger._log {
		log.Fatalf(format, v...)
	}
}
