package logger

import (
	"go-web-template/pkg/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var globalLogger *zap.Logger

var sugarLogger *zap.SugaredLogger

var isDev bool

func Init() {
	switch conf.AppConf.Mode {
	case "", "dev", "test", "debug", "local", "develop", "development":
		isDev = true
	}
	encoder := getEncoder()
	writeSyncer := getWriterSyncer()
	level := getLogLevel()
	core := zapcore.NewCore(encoder, writeSyncer, level)
	globalLogger = zap.New(core, zap.AddCaller())
	sugarLogger = globalLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	if isDev {
		devConf := zap.NewDevelopmentEncoderConfig()
		devConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(devConf)
	} else {
		prodConf := zap.NewProductionEncoderConfig()
		prodConf.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewJSONEncoder(prodConf)
	}
}

func getWriterSyncer() zapcore.WriteSyncer {
	if isDev {
		return zapcore.AddSync(os.Stdout)
	} else {
		lumberJackLogger := lumberjack.Logger{
			Filename:   conf.LogConf.FileName,
			MaxSize:    conf.LogConf.MaxSize,
			MaxAge:     conf.LogConf.MaxAge,
			MaxBackups: conf.LogConf.MaxBackups,
			Compress:   conf.LogConf.Compress,
		}
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberJackLogger))
	}
}

func getLogLevel() zapcore.Level {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(conf.LogConf.Level)); err != nil {
		panic(err)
	}
	return level
}
