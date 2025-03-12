package zapLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/**
 * @Description
 * @Date 2025/3/10 15:11
 **/

/*
	使用zap 日志库
   1. 该zap 日志库实现了加前缀
   2. 实现了按天进行分割日志
   3. zap.L() // 返回 *zap.Logger - 强类型、高性能日志接  zap.S() // 返回 *zap.SugaredLogger - 弱类型、开发友好型接口
// 使用 zap.L() 的典型场景（结构化日志）
zap.L().Debug("调试信息",
    zap.String("module", "network"),
    zap.Int("connections", 42))

// 使用 zap.S() 的典型场景（开发便捷）
zap.S().Debugw("连接状态",
    "module", "network",
    "active_conn", 42,
    "latency", 15.6)

zap.S().Debugf("当前内存使用: %.2f MB", 123.45)

*/

func InitZapLog() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 这里将输出的时间key 由原来的ts 改为 time
	cfg.EncoderConfig.TimeKey = "time"
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	writer := &dynamicLogWriter{
		logDir: "logs",
	}
	// 提前创建日志目录
	if err := os.MkdirAll(writer.logDir, 0755); err != nil {
		panic(err)
	}
	// 确保初始文件存在
	_, _ = writer.Write([]byte{})

	//
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(writer),
	)
	encoder := &prefixedEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
	}

	core := zapcore.NewCore(
		// 注意NewConsoleEncode 是控制台输出，控制台输出的话是文本格式，NewJSONEncoder 是json格式输出
		//zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		encoder,
		zapcore.AddSync(writeSyncer),
		zapcore.DebugLevel,
	)
	// 这里把zap.AddCaller()加上去，这样就可以显示调用的文件名和行号了
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)

}

// 定义前缀
const logPrefix = "[MyApp] "

type prefixedEncoder struct {
	zapcore.Encoder
}

func (pe *prefixedEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := pe.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	// 在日志行的最前面添加前缀
	logLine := buf.String()
	buf.Reset()
	buf.AppendString(logPrefix + logLine)
	return buf, nil

}

// 按天进行分割
type dynamicLogWriter struct {
	mu         sync.Mutex
	currentDay string
	file       *os.File
	logDir     string
}

func (w *dynamicLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 获取当前日期
	currentDay := time.Now().Format("2006-01-02")
	// 检查日期是否发生变化
	if currentDay != w.currentDay {
		if w.file != nil {
			if err := w.file.Close(); err != nil {
				return 0, fmt.Errorf("关闭日志文件失败: %w", err)
			}

		}

		filePath := filepath.Join(w.logDir, "app-"+currentDay+".log")
		file, er := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if er != nil {
			return 0, er
		}
		w.file = file
		w.currentDay = currentDay
	}
	// 如果日期没变的话就写入日志
	return w.file.Write(p)
}
func (w *dynamicLogWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}
