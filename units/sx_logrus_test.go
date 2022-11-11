package units

import (
	"sx.cases/sx_logrus"
	"testing"
)

var sxLogrusCfg = &sx_logrus.Cfg{
	Key:              "test.exp",
	Level:            6,
	LogFile:          "./logs/file.log",
	LogRotationTime:  0,
	LogMaxAge:        7200,
	LogRotationCount: 2,
}

// TestSXLogrus logrus日志库测试
func TestSXLogrus(t *testing.T) {
	sx_logrus.NewApi().New(sxLogrusCfg).GetClient().Infoln("info log")
	return
}
