package log

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YHemin/toolkit/log/hook"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	Path   string
	Level  string
	Fields map[string]string

	MaxAge        int
	DisableStdout bool
	Format        string // json/console/text
	RotateDay     int
}

// SugaredLogger ..
type SugaredLogger struct {
	*zap.SugaredLogger
}

// Logger ..
type Logger struct {
	*zap.Logger
	config *Config
}

var (
	// std is the name of the standard logger in stdlib `log`
	logger = &Logger{}
	sugger = &SugaredLogger{}
)

func init() {
	l, _ := zap.NewDevelopment()
	logger = &Logger{l, &Config{}}
	sugger = logger.Sugar()
}

// Sugar copy zaplog
func (log *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{log.Logger.Sugar()}
}

// New ..
func New(config *Config) (*Logger, error) {
	var (
		lvl        zapcore.Level
		err        error
		hooks      []zapcore.WriteSyncer
		rotatehook *rotatelogs.RotateLogs
		ecoder     zapcore.Encoder
		timeKey    = "time"
		levelKey   = "level"
		msgKey     = "msg"
	)
	if config.Level != "" {
		lvl = hook.ParseLevel(config.Level)
	}

	if config.Path != "" {
		dir := getDir(config.Path)
		if isPathNotExist(dir) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}

		if config.RotateDay != 0 {
			var fn = config.Path
			if !filepath.IsAbs(fn) {
				v, err := filepath.Abs(fn)
				if err != nil {
					return nil, err
				}
				fn = v
			}

			rotatehook, err = rotatelogs.New(
				fn+".%Y%m%d",
				rotatelogs.WithLinkName(fn),
				rotatelogs.WithMaxAge(time.Hour*24*time.Duration(config.MaxAge)),
				rotatelogs.WithRotationTime(time.Hour*24*time.Duration(config.RotateDay)),
			)
			hooks = append(hooks, zapcore.AddSync(rotatehook))
		}
	}

	if config.DisableStdout == false {
		hooks = append(hooks, os.Stdout)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        timeKey,
		LevelKey:       levelKey,
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     msgKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch strings.ToLower(config.Format) {
	case "json":
		ecoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		ecoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(
		ecoder,
		zapcore.NewMultiWriteSyncer(hooks...),
		lvl,
	))

	core := zapcore.NewTee(cores...)
	var l *zap.Logger
	l = zap.New(core)

	return &Logger{l, config}, nil
}

// Init ...
func Init(config *Config) error {
	var err error
	logger, err = New(config)
	if err != nil {
		return err
	}
	sugger = logger.Sugar()
	return nil
}

func isPathNotExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func getDir(path string) string {
	paths := strings.Split(path, "/")
	return strings.Join(
		paths[:len(paths)-1],
		"/",
	)
}
