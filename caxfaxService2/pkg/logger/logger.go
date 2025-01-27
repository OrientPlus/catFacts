package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	LogPath      string `yaml:"path"`
	Level        string `yaml:"level"`
	Name         string `yaml:"name"`
	ConsoleWrite bool   `yaml:"consolewrite"`
}

type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	Exit(code int)
	SetNoLock()
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer bytes.Buffer

	levelColors := map[logrus.Level]string{
		logrus.DebugLevel: "\033[36m", // Cyan
		logrus.InfoLevel:  "\033[32m", // Green
		logrus.WarnLevel:  "\033[33m", // Yellow
		logrus.ErrorLevel: "\033[31m", // Red
		logrus.FatalLevel: "\033[35m", // Magenta
		logrus.PanicLevel: "\033[41m", // Red Background
	}
	resetColor := "\033[0m"

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	levelColor, ok := levelColors[entry.Level]
	if !ok {
		levelColor = resetColor // По умолчанию
	}
	level := strings.ToUpper(entry.Level.String())

	file, line := findCaller()

	file = fmt.Sprintf("%s:%d", filepath.Base(file), line)

	logLine := fmt.Sprintf("%s %s<%s> %s%s %s\n",
		timestamp,
		levelColor, level,
		entry.Message, resetColor,
		file,
	)

	buffer.WriteString(logLine)
	return buffer.Bytes(), nil
}

func findCaller() (string, int) {
	const maxDepth = 25

	for i := 0; i < maxDepth; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if !strings.Contains(file, "/runtime/") && !strings.Contains(file, "logger.go") && !strings.Contains(file, "/logrus") {
			return filepath.Base(file), line
		}
	}

	return "unknown", 0
}

func NewLogger(config LoggerConfig) (Logger, error) {
	lg := logrus.Logger{}

	file, err := os.Create(config.LogPath + config.Name)
	if err != nil {
		return nil, err
	}

	var writer io.Writer
	if config.ConsoleWrite {
		writer = io.MultiWriter(os.Stdout, file)
	} else {
		writer = file
	}

	var level logrus.Level
	switch config.Level {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel

	}

	lg.SetOutput(writer)
	lg.SetLevel(level)
	lg.SetReportCaller(true)
	lg.SetFormatter(&CustomFormatter{})

	return &lg, nil
}
