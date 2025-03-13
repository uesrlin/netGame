package logrusLog

import (
	"bytes"
	"fmt"
	"github.com/pochard/logrotator"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

/**
 * @Description
 * @Date 2025/3/10 15:10
 **/

/*

使用方式
// 基础日志
 logrus.Info("服务启动成功")

// 带字段日志
 logrus.WithFields(logrus.Fields{
    "IP":    "192.168.1.1",
    "Port":  8080,
}).Warn("端口监听中")

// 格式化日志
 logrus.Debugf("当前队列长度: %d", 15)

logrus.Warn("端口监听中", "IP", "192.168.1.1", "Port", 8080) 高版本使用1.9+ 以上
*/

func InitLogrus() {
	const logDir = "server/logs"

	// 显式创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}
	rotator, err := logrotator.NewTimeBasedRotator(filepath.Join(logDir, "%Y%m%d-%H%M.log"), time.Hour*24)
	if err != nil {
		panic(fmt.Sprintf("logrotator.NewTimeBasedRotator Error %s exist", err))
	}

	writers := []io.Writer{
		rotator,
		os.Stdout}
	//  同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logrus.SetOutput(fileAndStdoutWriter)

	logrus.SetFormatter(&CustomFormatter{
		Prefix: "MYAPP",
		TextFormatter: logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel) //开启返回函数名和行号

}

// 前缀
type CustomFormatter struct {
	logrus.TextFormatter
	Prefix string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 创建缓冲区
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 设置颜色
	//var levelColor int
	//switch entry.Level {
	//case logrus.WarnLevel:
	//	levelColor = 33 // 黄色
	//case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
	//	levelColor = 31 // 红色
	//case logrus.DebugLevel:
	//	levelColor = 36 // 青色
	//default:
	//	levelColor = 34 // 蓝色
	//}
	// 构建基础日志信息
	fmt.Fprintf(b, "[%s]", f.Prefix)
	fmt.Fprintf(b, "  [%s]", entry.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(b, "  %-7s", strings.ToUpper(entry.Level.String())) // 带颜色的级别
	//fmt.Fprintf(b, "  \x1b[%dm[%s]\x1b[0m", levelColor, entry.Time.Format("2006-01-02 15:04:05"))
	//fmt.Fprintf(b, "  \x1b[%dm%-7s\x1b[0m", levelColor, strings.ToUpper(entry.Level.String())) // 带颜色的级别

	// 添加消息主体
	if entry.Message != "" {
		fmt.Fprintf(b, " - %s", entry.Message)
	}

	// 处理附加字段
	if len(entry.Data) > 0 {
		b.WriteString(" | ")
		fields := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		b.WriteString(strings.Join(fields, " "))
	}

	// 处理调用者信息
	if entry.HasCaller() {
		fmt.Fprintf(b, " %s:%d",
			path.Base(entry.Caller.File), // 仅显示文件名
			entry.Caller.Line)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
