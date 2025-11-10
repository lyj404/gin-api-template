package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogHook struct {
	infoFile *os.File // 普通日志文件
	errFile  *os.File // 错误日志文件
	fileDate string   // 日志文件日期
	logPath  string   // 日志文件路径
	infoSync zapcore.WriteSyncer
	errSync  zapcore.WriteSyncer
	mu       sync.Mutex
}

// NewZapLogHook 创建一个新的zap日志钩子
func NewZapLogHook(logPath string) (*ZapLogHook, error) {
	hook := &ZapLogHook{
		logPath: logPath,
	}

	// 初始化日志文件
	if err := hook.rotateFiles(time.Now().Format("2006-01-02")); err != nil {
		return nil, fmt.Errorf("failed to initialize log hook: %v", err)
	}

	return hook, nil
}

// rotateFiles 用于轮换日志文件
func (hook *ZapLogHook) rotateFiles(timer string) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 关闭当前日志
	if hook.infoFile != nil {
		hook.infoFile.Close()
	}
	if hook.errFile != nil {
		hook.errFile.Close()
	}

	// 创建日志目录
	dirName := filepath.Join(hook.logPath, timer)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	var err error

	// 创建普通日志文件
	infoFileName := filepath.Join(dirName, "info.log")
	hook.infoFile, err = os.OpenFile(infoFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create info log file: %v", err)
	}

	// 创建错误日志文件
	errFileName := filepath.Join(dirName, "error.log")
	hook.errFile, err = os.OpenFile(errFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create error log file: %v", err)
	}

	// 创建 lumberjack logger
	infoLumberjackLogger := &lumberjack.Logger{
		Filename:   infoFileName,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}
	errLumberjackLogger := &lumberjack.Logger{
		Filename:   errFileName,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}

	hook.infoSync = zapcore.AddSync(infoLumberjackLogger)
	hook.errSync = zapcore.AddSync(errLumberjackLogger)

	hook.fileDate = timer
	return nil
}

// GetInfoSync 获取 info 级别日志的 WriteSyncer
func (hook *ZapLogHook) GetInfoSync() zapcore.WriteSyncer {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 检查日期是否需要轮换
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != hook.fileDate {
		hook.rotateFiles(currentDate)
	}

	return hook.infoSync
}

// GetErrorSync 获取 error 级别日志的 WriteSyncer
func (hook *ZapLogHook) GetErrorSync() zapcore.WriteSyncer {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 检查日期是否需要轮换
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != hook.fileDate {
		hook.rotateFiles(currentDate)
	}

	return hook.errSync
}

// InitZapLogger 初始化并配置 zap.Logger
func InitZapLogger() *zap.Logger {
	// 创建日志钩子
	hook, err := NewZapLogHook("./logs")
	if err != nil {
		zap.L().Fatal("Failed to initialize log hook", zap.Error(err))
	}

	// 设置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建核心
	core := zapcore.NewTee(
		// 控制台输出
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		// info 文件输出
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), hook.GetInfoSync(), zapcore.InfoLevel),
		// error 文件输出
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), hook.GetErrorSync(), zapcore.ErrorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
