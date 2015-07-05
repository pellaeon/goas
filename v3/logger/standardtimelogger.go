// Copyright (C) 2015 Pellaeon Lin <pellaeon@cnmc.tw>
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger

import (
	"io"
	"sync"
	"time"
)

// StandardTimeLogger is a simple logger writing to the given writer with
// given time format. Beside custom time format, it is the same as
// StandardLogger
type StandardTimeLogger struct {
	mutex      sync.Mutex
	out        io.Writer
	timeformat string
}

// NewStandardTimeLogger creates the standard logger with custom time format.
func NewStandardTimeLogger(out io.Writer, tf string) Logger {
	return &StandardTimeLogger{out: out, timeformat: tf}
}

// Debug is specified on the Logger interface.
func (sl *StandardTimeLogger) Debug(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [DEBUG] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Info is specified on the Logger interface.
func (sl *StandardTimeLogger) Info(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [INFO] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Warning is specified on the Logger interface.
func (sl *StandardTimeLogger) Warning(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [WARNING] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Error is specified on the Logger interface.
func (sl *StandardTimeLogger) Error(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [ERROR] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Critical is specified on the Logger interface.
func (sl *StandardTimeLogger) Critical(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [CRITICAL] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Fatal is specified on the Logger interface.
func (sl *StandardTimeLogger) Fatal(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, time.Now().Format(sl.timeformat))
	io.WriteString(sl.out, " [FATAL] ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}
