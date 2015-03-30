package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/sentry"
	"os"
	"time"
)

var log *logrus.Logger
var dsn = "http://46ac2dc9fe4b48c38c277421b77dde52:19fe10508ad7402ebe0457970e92bc0e@128.199.210.115/2"

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
