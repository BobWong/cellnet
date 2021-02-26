package log

import "github.com/bobwong89757/golog/logs"

var mLog *logs.BeeLogger

func SetLog(v *logs.BeeLogger) {
	mLog = v
}

func GetLog() *logs.BeeLogger{
	return mLog
}
