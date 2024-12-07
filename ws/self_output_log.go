package ws

import (
	"log"
	"os"
)

// CreateLogger 创建一个新的log.Logger实例，并将输出重定向到指定的文件中
func CreateLogger(logFileName string) (*log.Logger, error) {
	// 打开或创建日志文件
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// 创建一个新的log.Logger实例
	logger := log.New(file, "", log.LstdFlags)

	return logger, nil
}

// 使用示例
/*  func main() {
	// 创建一个日志记录器，输出到"custom_log.log"文件
	logger, err := CreateLogger("custom_log.log")
	if err != nil {
		log.Fatalf("无法创建日志记录器: %v", err)
	}

	// 使用新的logger实例记录日志
	logger.Println("这是一条自定义日志文件名中的日志消息")
	logger.Printf("这是一条带有格式化的日志消息: %d", 42)

	// ... 你的其他代码 ...
}
 */