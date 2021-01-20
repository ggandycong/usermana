package log

import (
	"fmt"
	"log"
	"os"
)

// Logger level 日志等级 logPath日志文件.
type Logger struct {
	level   int
	logPath string
}

// 日志等级.
const (
	LevelFatal   int = iota // 严重错误信息.
	LevelError              // 错误信息.
	LevelWarning            // 警告信息.
	LevelInfo               // 普通信息.
	LevelDebug              // 调试信息.
)

var logger Logger

// Config 加载日志配置.
func Config(logPath string, level int) error {
	// 打开日志文件.
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	logger.logPath = logPath
	logger.level = level
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logFile)

	return nil
}

// Debugf 输出调试信息.
func Debugf(format string, v ...interface{}) {
	if logger.level >= LevelDebug {
		log.SetPrefix("debug ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Infof 输出普通信息.
func Infof(format string, v ...interface{}) {
	if logger.level >= LevelInfo {
		log.SetPrefix("info ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warningf 输出警告信息.
func Warningf(format string, v ...interface{}) {
	if logger.level >= LevelWarning {
		log.SetPrefix("warning ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Errorf 输出错误信息.
func Errorf(format string, v ...interface{}) {
	if logger.level >= LevelError {
		log.SetPrefix("error ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatalf 输出严重错误信息.
func Fatalf(format string, v ...interface{}) {
	if logger.level >= LevelFatal {
		log.SetPrefix("fatal ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}
