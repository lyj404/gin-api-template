package logger

import (
	"bytes"
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
)

// LogFormatter 自定义日志格式
type LogFormatter struct{}

// Format 实现 logrus.Formatter 接口，用于自定义日志格式
func (LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	// 根据日志级别设置颜色
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 37 // 灰色
	case logrus.WarnLevel:
		levelColor = 33 // 黄色
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // 红色
	default:
		levelColor = 36 // 蓝色
	}

	// 初始化用于存储日志的 buffer
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 格式化日志时间
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 格式化日志字段
	var fields string
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			fields += fmt.Sprintf("[%s: %v] ", key, value)
		}
	}

	// 如果日志包含调用者信息，则格式化并添加到日志中
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		fmt.Fprintf(b, "\x1b[%dm[%s] [%s] %s %s %s\x1b[0m\n", levelColor, timestamp, entry.Level, funcVal, fileVal, fields)
	} else {
		fmt.Fprintf(b, "\x1b[%dm[%s] [%s] %s\x1b[0m\n", levelColor, timestamp, entry.Level, fields)
	}

	return b.Bytes(), nil
}
