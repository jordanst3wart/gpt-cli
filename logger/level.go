package logger

import "github.com/rs/zerolog"

var levelMap = map[string]zerolog.Level{
	"nop":   zerolog.NoLevel,
	"error": zerolog.ErrorLevel,
	"warn":  zerolog.WarnLevel,
	"info":  zerolog.InfoLevel,
	"debug": zerolog.DebugLevel,
	"trace": zerolog.TraceLevel,
}

func LevelList() []string {
	return []string{"nop", "error", "warn", "info", "debug", "trace"}
}

func LevelFrom(s string) zerolog.Level {
	if lvl, ok := levelMap[s]; ok {
		return lvl
	}
	return zerolog.NoLevel
}
