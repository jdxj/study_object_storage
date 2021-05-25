package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	sugar *zap.SugaredLogger
)

func Init(fileName, appName string, maxSize, maxAge, maxBackups, level int, localTime, compress bool) {
	hook := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "file",
		FunctionKey:      "func",
		StacktraceKey:    "stack",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "",
	}

	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(level))

	writers := []zapcore.WriteSyncer{
		zapcore.AddSync(hook),
	}
	if zapcore.Level(level) == zapcore.DebugLevel {
		writers = append(writers, os.Stdout)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writers...),
		atomicLevel,
	)

	// options
	caller := zap.AddCaller()
	dev := zap.Development()
	field := zap.Fields(zap.String("appName", appName))

	logger := zap.New(core, caller, dev, field)
	sugar = logger.Sugar()
}

func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func DPanicf(format string, args ...interface{}) {
	sugar.DPanicf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	sugar.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}
