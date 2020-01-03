package main

import (
	"fmt"
	"go-snapshot/common"
	"os/exec"
	"sync"
)

//Logio info
func Logio(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "iostat")
	cmd := exec.Command("bash", "-c", "iostat 1 5")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

//Logmpstat info
func Logmpstat(childDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := fmt.Sprintf("%s/%s", childDir, "mpstat")
	cmd := exec.Command("bash", "-c", "iostat 1 5")
	out, _ := cmd.CombinedOutput()
	_ = common.CreateFileWriteNote(fileName, string(out))
}

