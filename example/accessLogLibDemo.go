package example

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	logger "github.com/wgffgw/loggerX"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*-----------------------------------logrusDemo-------------------------------------*/

func logrusTest() {
	//获取进程全局logrus实例
	log := NewLogrusLogger()
	//支持logrus日志库使用LoggerX接口格式打印日志
	logger.AccessLogLib("logrus", log)

	tag := "Tag message"
	err := fmt.Errorf("Logger Test")

	logger.Ix(context.Background(), tag, "this is info log:%v", err)
	logger.Wx(context.Background(), tag, "this is warning log:%v", err)
	logger.Ex(context.Background(), tag, "this is error log:%v", err)
	logger.Dx(context.Background(), tag, "this is debug log:%v", err)
	logger.Tx(context.Background(), tag, "this is trace log:%v", err)
	logger.F(tag, "this is trace log:%v", err)

}

func NewLogrusLogger() *logrus.Logger {
	//初始化全局logger实例
	LogrusT := logrus.New()
	//添加logger配置项
	LogrusT.SetLevel(logrus.InfoLevel)
	LogrusT.SetFormatter(&logrus.JSONFormatter{})
	//添加自定义logrus Hook接口
	LogrusT.AddHook(&DefaultFieldHook{})
	return LogrusT
}

//自定义Hook 接口
type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

/*-----------------------------------zapDemo-------------------------------------*/

func zapTest() {
	log := InitZap()
	defer log.Sync()
	//支持Zap日志库使用LoggerX接口格式打印日志
	logger.AccessLogLib("zap", log)

	tag := "Tag message"
	err := fmt.Errorf("Logger Test")

	logger.Ix(context.Background(), tag, "this is info log:%v", err)
	logger.Wx(context.Background(), tag, "this is warning log:%v", err)
	logger.Ex(context.Background(), tag, "this is error log:%v", err)
	logger.Dx(context.Background(), tag, "this is debug log:%v", err)
	//T -- zap debug
	logger.Tx(context.Background(), tag, "this is trace log:%v", err)
	logger.F(tag, "this is trace log:%v", err)
}

func InitZap() *zap.Logger {
	// 动态日志等级
	dynamicLevel := zap.NewAtomicLevel()

	core := zapcore.NewTee(
		// 有好的格式、输出控制台、动态等级
		zapcore.NewCore(zapcore.NewConsoleEncoder(NewEncoderConfig()), os.Stdout, dynamicLevel),
	)
	logger := zap.New(core, zap.AddCaller())

	// 将当前日志等级设置为Debug
	// 注意日志等级低于设置的等级，日志文件也不分记录
	dynamicLevel.SetLevel(zap.DebugLevel)
	return logger
}

// 格式化时间
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",                           // json时时间键
		LevelKey:       "L",                           // json时日志等级键
		NameKey:        "N",                           // json时日志记录器名
		CallerKey:      "C",                           // json时日志文件信息键
		MessageKey:     "M",                           // json时日志消息键
		StacktraceKey:  "S",                           // json时堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 友好日志等级名大小写（info INFO）
		EncodeTime:     TimeEncoder,                   // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 日志文件信息（包/文件.go:行号）
	}
}
