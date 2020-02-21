package main

import (
	"fmt"
	"go-snapshot/common"
	"os/exec"
	"sync"
	"time"
)

//LogIo info
func LogIo(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "iostat")
	cmd := exec.Command("bash", "-c", "iostat -m -x 1 5")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogMpstat info
func LogMpstat(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "mpstat")
	cmd := exec.Command("bash", "-c", "mpstat 1 5")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogDiskSpace info
func LogDiskSpace(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "disk_space")
	cmd := exec.Command("bash", "-c", "df -h")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogTop info
func LogTop(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "top")
	cmd := exec.Command("bash", "-c", "top -bn5")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogMemoInfo info
func LogMemoInfo(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "meminfo")
	cmd := exec.Command("bash", "-c", "cat /proc/meminfo")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogInterrupts info
func LogInterrupts(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "interrupts")
	cmd := exec.Command("bash", "-c", "cat /proc/interrupts")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogPs info
func LogPs(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "ps")
	cmd := exec.Command("bash", "-c", "ps -eaF")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogNetStat info
func LogNetStat(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "netstat")
	cmd := exec.Command("bash", "-c", "netstat -antp")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//LogInnodbStatus info
func LogInnodbStatus(conf *common.Config, childDir string, wg *sync.WaitGroup) {
	fmt.Println(time.Now().Format("2006-01-02-15-04-05"))
	db, err := common.NewMySQLConnection(conf)
	if err != nil {
		Log.Error("无法建立数据库连接，错误信息：%s", err)
		return
	}
	fmt.Println(time.Now().Format("2006-01-02-15-04-05"))
	defer func() { _ = db.Close() }()
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "innodb_status")
	innodbStaus, err := GetInnodbStaus(db)
	if err != nil {
		Log.Error("get innodb status failed")
	}
	_ = common.CreateFileWriteNote(fileName, innodbStaus)
}

//LogProcesslist info
func LogProcesslist(conf *common.Config, childDir string, wg *sync.WaitGroup) {
	db, err := common.NewMySQLConnection(conf)
	if err != nil {
		Log.Error("无法建立数据库连接，错误信息：%s", err)
		return
	}
	defer func() { _ = db.Close() }()
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "processlist")
	processList, err := GetProcesslist(db)
	if err != nil {
		Log.Error("get processlist failed")
	}
	_ = common.CreateFileWriteNote(fileName, processList)
}
