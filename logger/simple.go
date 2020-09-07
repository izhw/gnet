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
	"sync/atomic"
)

var _ Logger = &simpleLogger{}

func NewSimpleLogger() Logger {
	return &simpleLogger{
		depth:  2,
		level:  int32(InfoLevel),
		logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}
}

func NewSimpleLoggerWithLevel(l Level) Logger {
	return &simpleLogger{
		depth:  2,
		level:  int32(l),
		logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}
}

// simpleLogger implementation for Logger
type simpleLogger struct {
	depth  int
	level  int32
	logger *log.Logger
}

func (l *simpleLogger) Debug(v ...interface{}) {
	if l.getLevel() > DebugLevel {
		return
	}
	l.logger.Output(l.depth, decorate(DebugLevel, fmt.Sprintln(v...)))
}

func (l *simpleLogger) Debugf(format string, v ...interface{}) {
	if l.getLevel() > DebugLevel {
		return
	}
	l.logger.Output(l.depth, decorate(DebugLevel, fmt.Sprintf(format, v...)))
}

func (l *simpleLogger) Info(v ...interface{}) {
	if l.getLevel() > InfoLevel {
		return
	}
	l.logger.Output(l.depth, decorate(InfoLevel, fmt.Sprintln(v...)))
}

func (l *simpleLogger) Infof(format string, v ...interface{}) {
	if l.getLevel() > InfoLevel {
		return
	}
	l.logger.Output(l.depth, decorate(InfoLevel, fmt.Sprintf(format, v...)))
}

func (l *simpleLogger) Warn(v ...interface{}) {
	if l.getLevel() > WarnLevel {
		return
	}
	l.logger.Output(l.depth, decorate(WarnLevel, fmt.Sprintln(v...)))
}

func (l *simpleLogger) Warnf(format string, v ...interface{}) {
	if l.getLevel() > WarnLevel {
		return
	}
	l.logger.Output(l.depth, decorate(WarnLevel, fmt.Sprintf(format, v...)))
}

func (l *simpleLogger) Error(v ...interface{}) {
	if l.getLevel() > ErrorLevel {
		return
	}
	l.logger.Output(l.depth, decorate(ErrorLevel, fmt.Sprintln(v...)))
}

func (l *simpleLogger) Errorf(format string, v ...interface{}) {
	if l.getLevel() > ErrorLevel {
		return
	}
	l.logger.Output(l.depth, decorate(ErrorLevel, fmt.Sprintf(format, v...)))
}

func (l *simpleLogger) Fatal(v ...interface{}) {
	if l.getLevel() > FatalLevel {
		return
	}
	l.logger.Output(l.depth, decorate(FatalLevel, fmt.Sprintln(v...)))
	os.Exit(1)
}

func (l *simpleLogger) Fatalf(format string, v ...interface{}) {
	if l.getLevel() > FatalLevel {
		return
	}
	l.logger.Output(l.depth, decorate(FatalLevel, fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *simpleLogger) GetLevel() Level {
	return l.getLevel()
}

func (l *simpleLogger) SetLevel(level Level) {
	atomic.StoreInt32(&l.level, int32(level))
}

func (l *simpleLogger) getLevel() Level {
	return Level(atomic.LoadInt32(&l.level))
}

func decorate(logLevel Level, msg string) string {
	return fmt.Sprintf("[%s] %s", logLevel, msg)
}
