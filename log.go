package logger

import (
	"fmt"
	"io"
	"logger/internal/libs/parsers/json"
	"logger/internal/libs/parsers/text"
	"os"
	"strings"
)

type Logger struct {
	writer io.Writer
	level  Level
	Parser
}

type Parser interface {
	Parse(attrs map[string]interface{}) ([]byte, error)
}

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var Log = New()

func New() *Logger {
	logConfig, err := loadConfig()
	newLogger := Logger{writer: os.Stdout, level: InfoLevel, Parser: json.New()}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fallback to default, failed load config file: %s\n", err)
		return &newLogger
	}
	if logConfig.Output != "" {
		file, err := os.OpenFile(logConfig.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fallback to default, failed to open log file: %s\n", err)
			return &newLogger
		}
		newLogger.writer = file
	}
	newLogger.setLogLevel(logConfig.Level)
	newLogger.setParser(logConfig.Parser)

	return &newLogger
}

func Debug(traceId, message string, attributes map[string]interface{}) {
	Log.handleLine(traceId, message, attributes, DebugLevel)
}

func Info(traceId, message string, attributes map[string]interface{}) {
	Log.handleLine(traceId, message, attributes, InfoLevel)
}

func Warn(traceId, message string, attributes map[string]interface{}) {
	Log.handleLine(traceId, message, attributes, WarnLevel)
}

func Error(traceId, message string, attributes map[string]interface{}) {
	Log.handleLine(traceId, message, attributes, ErrorLevel)
}
func (l *Logger) handleLine(traceId, message string, attributes map[string]interface{}, level Level) {
	logIt := Log.canLog(level)
	if !logIt {
		return
	}
	attributes["level"] = getLevelAsString(level)
	attributes["traceID"] = traceId
	attributes["message"] = message
	line, err := l.Parse(attributes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse log: %s\n", err)
		return
	}
	err = Log.write(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: %s\n", err)
	}
}

func (l *Logger) write(line []byte) error {
	line = append(line, '\n')
	_, err := l.writer.Write(line)
	return err
}

func (l *Logger) canLog(level Level) bool {
	return l.level <= level
}

func (l *Logger) setLogLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		l.level = DebugLevel
	case "info":
		l.level = InfoLevel
	case "warn":
		l.level = WarnLevel
	case "error":
		l.level = ErrorLevel
	default:
		l.level = InfoLevel
	}
}

func getLevelAsString(level Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

func (l *Logger) setParser(parser string) {
	if strings.EqualFold(parser, "text") {
		l.Parser = text.New()
		return
	}
	l.Parser = json.New()
}
