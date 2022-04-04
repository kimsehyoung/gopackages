package shlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger   *logrus.Logger
	logLevel = logrus.DebugLevel
)

const (
	timestampFormat  = "2006-01-02 15:04:05" // refer to time.RFC3339
	maxMessageLength = "256"
	logCallStack     = "shlog.Logf" // package.function
	maxCallerDepth   = 15
)
const (
	FG_BLACK = 30 + iota
	FG_RED
	FG_GREEN
	FG_YELLOW
	FG_BLUE
	FG_MAGENTA
	FG_CYAN
	FG_WHITE
)

// Formatter implements logrus.Formatter interface.
type CustomFormatter struct {
	TimestampFormat  string
	MaxMessageLength string
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
}

// Default is stdout, DebugLevel.
// Use ChangeLogLevel to change log level.
func InitLogger(filename string) {
	logger = logrus.New()

	if filename != "" {
		out := &lumberjack.Logger{
			Filename: filename,
			MaxSize:  1, // MB
			// MaxAge: 10, // days
			// MaxBackups: 3,
			// LocalTime: ,
			// Compress: true,
		}
		logger.SetOutput(out)
	} else {
		logger.SetOutput(os.Stdout)
	}
	logger.SetReportCaller(true)
	logger.SetFormatter(&CustomFormatter{
		TimestampFormat:  timestampFormat,
		MaxMessageLength: maxMessageLength,
		CallerPrettyfier: callerPrettyfier,
	})
	logger.SetLevel(logLevel)
}

func callerPrettyfier(f *runtime.Frame) (function string, file string) {
	pc := make([]uintptr, maxCallerDepth)
	_ = runtime.Callers(0, pc)

	frames := runtime.CallersFrames(pc)
	for {
		frame, _ := frames.Next()

		if strings.Contains(frame.Function, logCallStack) {
			frame, _ = frames.Next()
			function := frame.Function[strings.LastIndex(frame.Function, ".")+1:]
			file := filepath.Base(filepath.Dir(frame.File))
			return function, file
		}
	}
}

// Use this after InitLogger
// level: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC
func ChangeLogLevel(level string) {
	switch level {
	case "PANIC":
		logger.SetLevel(logrus.PanicLevel)
	case "FATAL":
		logger.SetLevel(logrus.FatalLevel)
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logger.SetLevel(logrus.WarnLevel)
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
	case "TRACE":
		logger.SetLevel(logrus.TraceLevel)
	case "INFO": // INFO
		logger.SetLevel(logrus.InfoLevel)
	default:
		fmt.Println("Invalid log level. Refer to annotion")
	}
}

// Format building log message.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// output buffer
	var b bytes.Buffer

	// time stamp
	b.WriteString(entry.Time.Format(f.TimestampFormat))

	// colored log level
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = FG_GREEN
	case logrus.WarnLevel:
		levelColor = FG_YELLOW
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = FG_RED
	default:
		levelColor = FG_CYAN
	}
	b.WriteString(fmt.Sprintf(" [\x1b[%dm%s\x1b[0m]", levelColor, strings.ToUpper(entry.Level.String())))

	// file, path
	function, file := f.CallerPrettyfier(entry.Caller)
	b.WriteString(fmt.Sprintf("[%s][%s]", file, function))

	// message
	message := " %." + f.MaxMessageLength + "s"
	b.WriteString(fmt.Sprintf(message, entry.Message))

	// fields
	for k, v := range entry.Data {
		b.WriteString(fmt.Sprintf(" \x1b[%dm%s\x1b[0m=", levelColor, k))
		b.WriteString(fmt.Sprintf("%v", v))
	}

	b.WriteString("\n")
	return b.Bytes(), nil
}

// level: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC
func Log(level string, args ...interface{}) {
	switch level {
	case "PANIC":
		logger.Panic(args...)
	case "FATAL":
		logger.Fatal(args...)
	case "ERROR":
		logger.Error(args...)
	case "WARN":
		logger.Warn(args...)
	case "DEBUG":
		logger.Debug(args...)
	case "TRACE":
		logger.Trace(args...)
	default: // INFO
		logger.Info(args...)
	}
}

// level: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC
func Logf(level string, format string, args ...interface{}) {
	switch level {
	case "PANIC":
		logger.Panicf(format, args...)
	case "FATAL":
		logger.Fatalf(format, args...)
	case "ERROR":
		logger.Errorf(format, args...)
	case "WARN":
		logger.Warnf(format, args...)
	case "DEBUG":
		logger.Debugf(format, args...)
	case "TRACE":
		logger.Tracef(format, args...)
	default: // INFO
		logger.Infof(format, args...)
	}
}

func LogFields(level string, message string, object interface{}) {

	var elmts reflect.Value
	if t := reflect.TypeOf(object); t.Kind() == reflect.Ptr {
		elmts = reflect.ValueOf(object).Elem()
	} else {
		elmts = reflect.ValueOf(object)
	}

	fields := make(map[string]interface{})

	for i := 0; i < elmts.NumField(); i++ {
		field := elmts.Type().Field(i).Name
		v := elmts.Field(i)
		fields[field] = v
	}

	switch level {
	case "PANIC":
		logger.WithFields(fields).Panic(message)
	case "FATAL":
		logger.WithFields(fields).Fatal(message)
	case "ERROR":
		logger.WithFields(fields).Error(message)
	case "WARN":
		logger.WithFields(fields).Warn(message)
	case "DEBUG":
		logger.WithFields(fields).Debug(message)
	case "TRACE":
		logger.WithFields(fields).Trace(message)
	default: // INFO
		logger.WithFields(fields).Info(message)

	}
}
