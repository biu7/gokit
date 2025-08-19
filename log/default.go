package log

var Default Logger = With(
	NewLogger(LevelInfo),
	"caller", Caller(5),
)

func Info(msg string, args ...any) {
	Default.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Default.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Default.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default.Error(msg, args...)
}
