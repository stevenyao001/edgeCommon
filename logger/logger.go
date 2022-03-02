package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var logger *zap.SugaredLogger

func ErrorLog(key, desc, traceId string, err interface{}) {
	logger.Errorw(key, "desc", desc, "trace_id", traceId, "err", err)
}

func WarnLog(key, desc, traceId string, err interface{}) {
	logger.Warnw(key, "desc", desc, "trace_id", traceId, "err", err)
}

func InfoLog(key, desc, traceId string, detail interface{}) {
	logger.Infow(key, "desc", desc, "trace_id", traceId, "detail", detail)
}

func DebugLog(key, desc, traceId string, detail interface{}) {
	logger.Debugw(key, "desc", desc, "trace_id", traceId, "detail", detail)
}

/**
 * 获取日志
 * Filename 日志文件路径
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * level 日志级别
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 */
func newWriter(filePath string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}
}

//公用编码器
func newEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "T",
		LevelKey:    "L",
		CallerKey:   "F",
		MessageKey:  "K",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder, //大写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func InitLog(logPath string) {

	//输出格式
	encoder := newEncoder()

	//打开一个文件
	writer := newWriter(logPath)

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.DebugLevel)

	//同时写文件/打印
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writer))

	//组装核心
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)

	log := zap.New(core, zap.AddCaller()).WithOptions(zap.AddCallerSkip(1))
	logger = log.Sugar()
}
