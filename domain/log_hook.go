package domain

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LogHook struct {
	file     *os.File // 普通日志文件
	errFile  *os.File // 错误日志文件
	fileDate string   // 日志文件日期
	logPath  string   // 日志文件路径
	mu       sync.Mutex
}

// Fire 实现logrus.Hook接口，用于处理日志
func (hook *LogHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 获取当前日期
	timer := entry.Time.Format("2006-01-02")
	// 如果日志文件日期与当前日期不同，则轮换日志文件
	if timer != hook.fileDate {
		if err := hook.rotateFiles(timer); err != nil {
			return err
		}
	}

	// 格式化日志内容
	line, err := entry.String()
	if err != nil {
		return fmt.Errorf("failed to format log entry: %v", err)
	}

	// 将日志写入到普通日志文件
	if _, err := hook.file.Write([]byte(line)); err != nil {
		return fmt.Errorf("failed to write log to file: %v", err)
	}

	// 将错误日志写入到错误日志文件
	if entry.Level <= logrus.ErrorLevel {
		if _, err := hook.errFile.Write([]byte(line)); err != nil {
			return fmt.Errorf("failed to write to error log file: %v", err)
		}
	}
	return nil
}

// Levels 实现 logrus.Hook 接口，指定需要处理的日志级别
func (hook *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// rotateFiles 用户轮换日志文件
func (hook *LogHook) rotateFiles(timer string) error {
	// 关闭当前日志
	if hook.file != nil {
		if err := hook.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %v", err)
		}
	}

	// 关闭错误日志
	if hook.errFile != nil {
		if err := hook.errFile.Close(); err != nil {
			return fmt.Errorf("failed to close error log file: %v", err)
		}
	}

	// 创建日志目录
	dirName := fmt.Sprintf("%s/%s", hook.logPath, timer)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	var err error

	// 创建普通日志文件
	infoFileName := fmt.Sprintf("%s/info.log", dirName)
	hook.file, err = os.OpenFile(infoFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	// 创建错误日志文件
	errFileName := fmt.Sprintf("%s/error.log", dirName)
	hook.errFile, err = os.OpenFile(errFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create error log file: %v", err)
	}

	hook.fileDate = timer
	return nil
}

// InitLogger 初始化并配置 logrus.Logger
func InitLogger() *logrus.Logger {
	// 创建 logrus 实例
	logger := logrus.New()

	// 设置自定义日志格式化
	logger.SetFormatter(&LogFormatter{})

	// 初始化日志钩子
	hook := &LogHook{
		logPath: "./logs", // 设置日志文件存储路径
	}
	// 初始化日志文件
	if err := hook.rotateFiles(time.Now().Format("2006-01-02")); err != nil {
		logger.Fatal("Failed to initialize log hook:", err)
	}
	// 将钩子添加到 logrus 中
	logger.AddHook(hook)

	// 设置 logger 输出调用者信息（文件名和行号）
	logger.SetReportCaller(true)

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
