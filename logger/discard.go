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

import "os"

var _ Logger = &discardLogger{}

func NewDiscardLogger() Logger {
	return &discardLogger{}
}

// discardLogger implementation for Logger
type discardLogger struct{}

func (*discardLogger) Debug(...interface{}) {}

func (*discardLogger) Debugf(string, ...interface{}) {}

func (*discardLogger) Info(...interface{}) {}

func (*discardLogger) Infof(string, ...interface{}) {}

func (*discardLogger) Warn(...interface{}) {}

func (*discardLogger) Warnf(string, ...interface{}) {}

func (*discardLogger) Error(...interface{}) {}

func (*discardLogger) Errorf(string, ...interface{}) {}

func (*discardLogger) Fatal(...interface{}) {
	os.Exit(1)
}

func (*discardLogger) Fatalf(string, ...interface{}) {
	os.Exit(1)
}
