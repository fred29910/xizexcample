package logger

import (
	"io"
	"log"
	"os"
)

var (
	// InfoLogger 用于记录参考信息
	InfoLogger *log.Logger
	// ErrorLogger 用于记录错误信息
	ErrorLogger *log.Logger
)

func init() {
	// 默认初始化，将日志输出到标准输出和标准错误
	// 之后可以被 T041 中的配置管理覆盖
	InitLogger(os.Stdout, os.Stderr)
}

// InitLogger 初始化日志记录器，允许将输出重定向到不同的 io.Writer
func InitLogger(infoHandle io.Writer, errorHandle io.Writer) {
	InfoLogger = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
