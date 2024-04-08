package logger

import (
	"github.com/op/go-logging"
	"os"
	"strings"
)

var Log *logging.Logger

func init() {
	var log = logging.MustGetLogger("initialization")
	var level logging.Level
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{color}%{level:.6s} ▶%{color:reset} %{message}`,
	)

	var prebackend = logging.NewLogBackend(os.Stdout, "", 0)
	var formattedBackend = logging.NewBackendFormatter(prebackend, format)
	var backend = logging.AddModuleLevel(formattedBackend)

	logLevelEnv, ok := os.LookupEnv("LOG_LEVEL")

	if !ok || logLevelEnv == "" {
		log.Info("No LOG_LEVEL environment variable found, defaulting to INFO")
		level = logging.INFO
	} else {
		level, ok = ParseStringToLevel(logLevelEnv)
		if !ok {
			log.Errorf("Invalid log level: %s  Defaulting to INFO", logLevelEnv)
			level = logging.INFO
		} else {
			log.Infof("Setting log level to: %s", strings.ToUpper(logLevelEnv))
		}
	}

	backend.SetLevel(level, "mmock")
	logging.SetBackend(backend)

	Log = logging.MustGetLogger("mmock")
	log.Infof("Logger initialized with level: %s", strings.ToUpper(logging.GetLevel("mmock").String()))
}

var levelMap = map[string]logging.Level{
	"CRITICAL": logging.CRITICAL,
	"ERROR":    logging.ERROR,
	"WARNING":  logging.WARNING,
	"NOTICE":   logging.NOTICE,
	"INFO":     logging.INFO,
	"DEBUG":    logging.DEBUG,
}

func ParseStringToLevel(str string) (logging.Level, bool) {
	level, ok := levelMap[strings.ToUpper(str)]
	return level, ok
}
