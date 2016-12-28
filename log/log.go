package log

import "os"

var (
	RootContext *LogContext = NewLogContext()
	logger      Logger      = NewStdLogger(KVFormatter)
	level                   = TraceLevel
)

func SetLogger(l Logger) {
	logger = l
}

func GetLogger() Logger {
	return logger
}

func SetLevelByName(lvlName LevelName) {
	level = lvlName.ToLevel()
}

func GetLevel() Level {
	return level
}

func Enabled(lvl Level) bool {
	return lvl >= level
}

func Trace(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(TraceLevel) {
		logger.Log(ctx.LogContext().newEntry(TraceLevel, msg, keyvals))
	}
}

func Debug(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(DebugLevel) {
		logger.Log(ctx.LogContext().newEntry(DebugLevel, msg, keyvals))
	}
}

func Info(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(InfoLevel) {
		logger.Log(ctx.LogContext().newEntry(InfoLevel, msg, keyvals))
	}
}

func Warning(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(WarningLevel) {
		logger.Log(ctx.LogContext().newEntry(WarningLevel, msg, keyvals))
	}
}

func Error(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(ErrorLevel) {
		logger.Log(ctx.LogContext().newEntry(ErrorLevel, msg, keyvals))
	}
}

func Fatal(ctx Context, msg string, keyvals ...interface{}) {
	if Enabled(FatalLevel) {
		logger.Log(ctx.LogContext().newEntry(FatalLevel, msg, keyvals))
	}
	os.Exit(-1)
}
