package logger

import (
	"go.uber.org/zap/zapcore"
)

// LogFormatter 自定义日志格式
type LogFormatter struct{}

// Format 实现 zapcore.ObjectEncoder 接口，用于自定义日志格式
func (f *LogFormatter) Format(entry zapcore.Entry, enc zapcore.PrimitiveArrayEncoder) error {
	// 格式化日志时间
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 添加日志级别
	enc.AppendString("level")
	enc.AppendString(entry.Level.String())

	// 添加时间戳
	enc.AppendString("timestamp")
	enc.AppendString(timestamp)

	// 添加日志消息
	enc.AppendString("message")
	enc.AppendString(entry.Message)

	// 添加调用者信息
	if entry.Caller.Defined {
		enc.AppendString("caller")
		enc.AppendString(entry.Caller.String())
	}

	return nil
}

// NewJSONEncoder 创建自定义JSON编码器
func NewJSONEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}
