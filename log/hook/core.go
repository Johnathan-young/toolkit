package hook

import (
	"go.uber.org/zap/zapcore"
)

// ParseLevel .. parse level
func ParseLevel(loglevel string) zapcore.Level {
	var lv zapcore.Level
	lv.UnmarshalText([]byte(loglevel))
	return lv
}
