package common

import (
	"fmt"
	"os"
)

//CreateFileWriteNote for create stat file and write stat
func CreateFileWriteNote(fileName string, note string) (err error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.WriteString(note)
	return
}

//CompatibleLog for making log_file and log_dir compatible
func CompatibleLog(conf *Config) string {
	logDir := conf.Base.LogDir
	logFile := conf.Base.LogFile
	initLogFile := "snapshot.log"
	if logFile != "" {
		return logFile
	}
	if logDir != "" {
		return fmt.Sprintf("%s/%s", logDir, initLogFile)
	}
	return initLogFile
}
