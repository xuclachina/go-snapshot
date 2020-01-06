package common

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

// MyNewLogger constructor of needed logger
func MyNewLogger(conf *Config, logFile string) *logs.BeeLogger {
	return loggerInit(conf, logFile)
}

func loggerInit(conf *Config, logFile string) (log *logs.BeeLogger) {
	log = logs.NewLogger(0)
	log.EnableFuncCallDepth(true)
	log.SetLevel(conf.Base.LogLevel)
	if conf.Base.LogDir == "console" {
		_ = log.SetLogger("console")
	} else {
		_ = log.SetLogger("file", fmt.Sprintf(`{"filename":"%s", "level":%d, "maxlines":0,"maxsize":0, "daily":false, "maxdays":0}`, logFile, conf.Base.LogLevel))
	}
	return
}
