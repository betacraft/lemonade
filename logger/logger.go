package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/sentry"
	"os"
	"time"
)

var log *logrus.Logger
var dsn = "http://be0b18ca7977486ca1d93252015b01cb:1df514f1c72d49598d2d0fb4c1a67c1f@128.199.210.115/3"

func Init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)
	if os.Getenv("ENV") == "dev" {
		return
	}
	hook, err := logrus_sentry.NewSentryHook(dsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	hook.Timeout = 20 * time.Second
	if err != nil {
		logrus.Panic(err)
	}
	log.Debug("Adding sentry hook")
	log.Hooks.Add(hook)
}

func Get() *logrus.Logger {
	return log
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Panic(args ...interface{}) {
	log.Panic(args)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Err(args ...interface{}) {
	log.Error(args)
}
