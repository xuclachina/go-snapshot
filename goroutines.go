package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"go-snapshot/common"
	"os/exec"
	"sync"
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

//LogInnodbStatus info
func LogInnodbStatus(db mysql.Conn, childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "innodb_status")
	innodbStaus, err := GetInnodbStaus(db)
	if err != nil {
		Log.Error("get innodb status failed")
	}
	_ = common.CreateFileWriteNote(fileName, innodbStaus)
}
