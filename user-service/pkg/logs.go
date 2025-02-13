package pkg

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileHook is a custom hook for logging to a file with a different formatter
type FileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	LevelsVal []logrus.Level
}

func NewFileHook(levels []logrus.Level, writer io.Writer, formatter logrus.Formatter) *FileHook {
	return &FileHook{
		Writer:    writer,
		Formatter: formatter,
		LevelsVal: levels,
	}
}

func (hook *FileHook) Levels() []logrus.Level {
	return hook.LevelsVal
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	if os.Getenv("port") == "8080" {
		entry.Data["Environment"] = "Development"
	} else {
		entry.Data["Environment"] = "Production"
	}

	line, err := hook.Formatter.Format(entry)
	if err != nil {
		logrus.Errorf("Error formatting log entry for file: %v", err)
		return err
	}

	// Write the formatted entry to the writer
	_, err = hook.Writer.Write(line)
	if err != nil {
		logrus.Errorf("Error writing log entry to file: %v", err)
		return err
	}
	return nil
}

// SetupLogger initializes the logger with both terminal and file logging
func SetupLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatalf("Failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")

	fileLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     28, // Days
		Compress:   true,
	}

	// Set up terminal logger (this will be the default output)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		ForceColors:            true,  // Enable colors for terminal output
		DisableColors:          false, // Keep colors in terminal
		QuoteEmptyFields:       true,
		DisableQuote:           true,
		DisableLevelTruncation: true,
		PadLevelText:           false,
	})

	// Set log level
	logrus.SetLevel(logrus.InfoLevel)

	// Add custom hook for file logging with a different formatter (no colors)
	logrus.AddHook(NewFileHook(logrus.AllLevels, fileLogger, &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		ForceColors:      false, // Disable colors for file output
		DisableColors:    true,
		QuoteEmptyFields: true,
	}))
}

// LogrusLogger is a custom middleware for logging HTTP requests using logrus
func LogrusLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// save start time
		start := time.Now()
		err := next(c)
		latency := time.Since(start)
		req := c.Request()
		res := c.Response()
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableQuote:    true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			DisableColors:   false,
			PadLevelText:    false,
		})
		logrus.WithFields(logrus.Fields{
			"method":         req.Method,
			"uri":            req.RequestURI,
			"status":         res.Status,
			"latency":        latency,
			"ip":             c.RealIP(),
			"user_agent":     req.UserAgent(),
			"host":           req.Host,
			"referer":        req.Referer(),
			"protocol":       req.Proto,
			"content_length": req.ContentLength,
			"query_params":   c.QueryParams().Encode(),
			"response_size":  res.Size,
			"cookies":        req.Cookies(),
			"request_id":     c.Response().Header().Get(echo.HeaderXRequestID),
			"handler_name":   runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name(),
		}).Info("HTTP request")
		return err
	}
}
