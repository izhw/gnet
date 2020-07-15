// Copyright (c) 2020 izhw
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package logger

import (
	"fmt"
	"log"
	"os"
)

var _ Logger = &simpleLogger{}

func NewSimpleLogger() Logger {
	return &simpleLogger{
		depth:  2,
		logger: log.New(os.Stderr, "[network] ", log.LstdFlags|log.Lshortfile),
	}
}

// simpleLogger implementation for Logger
type simpleLogger struct {
	depth  int
	logger *log.Logger
}

func (l *simpleLogger) Debug(v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintln(v...))
}

func (l *simpleLogger) Debugf(format string, v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Info(v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintln(v...))
}

func (l *simpleLogger) Infof(format string, v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Warn(v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintln(v...))
}

func (l *simpleLogger) Warnf(format string, v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Error(v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintln(v...))
}

func (l *simpleLogger) Errorf(format string, v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintf(format, v...))
}

func (l *simpleLogger) Fatal(v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *simpleLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Output(l.depth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
