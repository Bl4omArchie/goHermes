package hermes

import (
	"log/slog"
	"os"
)


type SlogWrapper struct {
	Dir	string
	Log	*slog.Logger
}


func NewSlogWrapper(dir string, slogger *slog.Logger) *SlogWrapper {
	return &SlogWrapper {
		Dir: dir,
		Log: slogger,
	}
}


func InitializeSlogWrapper(dir string, logFormat string) *SlogWrapper {
	var logger *slog.Logger

	switch logFormat {
	case "text":
		f, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = slog.New(slog.NewTextHandler(f, nil))

	case "json":
		f, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = slog.New(slog.NewTextHandler(f, nil))
	}

	
	return NewSlogWrapper(dir, logger)
}
