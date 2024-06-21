/**
  @author: decision
  @date: 2024/6/14
  @note: 日志格式化
**/

package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"strings"
)

type CustomFormatter struct{}

// ANSI 颜色代码
const (
	AnsiReset     = 0
	AnsiRed       = 31
	AnsiHiRed     = 91
	AnsiGreen     = 32
	AnsiHiGreen   = 92
	AnsiYellow    = 33
	AnsiHiYellow  = 93
	AnsiBlue      = 34
	AnsiHiBlue    = 94
	AnsiMagenta   = 35
	AnsiHiMagenta = 95
	AnsiCyan      = 36
	AnsiHiCyan    = 96
	AnsiWhite     = 37
	AnsiHiWhite   = 97
)

// Format 自定义的 logrus 日志格式化器
// [Level - Datetime] content (file: line - function)
// Fields:
//   - "": ""
//   - "": ""
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var log, fields string

	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		log = fmt.Sprintf("[%5s - %s] %s (%s: %d - %s)\n", f.levelColorStr(entry.Level),
			timestamp, entry.Message, fName, entry.Caller.Line, entry.Caller.Function)
	} else {
		log = fmt.Sprintf("[%5s - %s] %s\n", f.levelColorStr(entry.Level),
			timestamp, entry.Message)
	}

	if len(entry.Data) > 0 {
		fields = "Fields:\n"
		for k, v := range entry.Data {
			fields += fmt.Sprintf("      - %s: %s\n", k, v)
		}

		log += fields
	}

	b.WriteString(log)
	return b.Bytes(), nil
}

func (f *CustomFormatter) levelColorStr(level logrus.Level) string {
	var levelColor int
	switch level {
	case logrus.DebugLevel:
		levelColor = AnsiCyan
	case logrus.InfoLevel:
		levelColor = AnsiGreen
	case logrus.WarnLevel:
		levelColor = AnsiYellow
	case logrus.ErrorLevel:
		levelColor = AnsiRed
	case logrus.FatalLevel:
		levelColor = AnsiMagenta
	case logrus.PanicLevel:
		levelColor = AnsiMagenta
	case logrus.TraceLevel:
		levelColor = AnsiBlue
	}

	return "\033[" + strconv.Itoa(levelColor) + "m" + strings.ToUpper(level.
		String()) + "\033[0m"
}
