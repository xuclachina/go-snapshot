package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/ziutek/mymysql/mysql"
	"os"
	"sync"
	"time"

	"go-snapshot/common"

	"github.com/astaxie/beego/logs"
)

//Log logger of project
var Log *logs.BeeLogger

func main() {
	// parse config file
	var confFile string
	var isnapshot bool

	flag.StringVar(&confFile, "c", "config.cfg", "snapshot configure file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(fmt.Sprintf("%10s: %s", "Version", "1.0.0"))
		os.Exit(0)
	}
	conf, err := common.NewConfig(confFile)
	if err != nil {
		fmt.Printf("NewConfig Error: %s\n", err.Error())
		return
	}
	if conf.Base.LogDir != "" {
		err = os.MkdirAll(conf.Base.LogDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll Error: %s\n", err.Error())
			return
		}
	}
	if conf.Base.SnapshotDir != "" {
		err = os.MkdirAll(conf.Base.SnapshotDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll Error: %s\n", err.Error())
			return
		}
	}

	// init log and other necessary
	Log = common.MyNewLogger(conf, common.CompatibleLog(conf))

	db, err := common.NewMySQLConnection(conf)
	if err != nil {
		fmt.Printf("NewMySQLConnection Error: %s\n", err.Error())
		return
	}
	defer func() { _ = db.Close() }()

	// start...
	Log.Info("start snapshot...")
	go timeout()
	for {
		isnapshot = checkCondition(conf, db)
		now := time.Now()
		childDir := fmt.Sprintf("%s/%s", conf.Base.SnapshotDir, now)
		if isnapshot {
			makeSnapshot(childDir)
		}
		time.Sleep(time.Second * Interval)
	}

}

func timeout() {
	time.AfterFunc(TimeOut*time.Second, func() {
		Log.Error("Execute timeout")
		os.Exit(1)
	})
}

func checkCondition(conf *common.Config, db mysql.Conn) (result bool) {
	metrics := make(map[string]int)
	c1, _ := cpu.Times(false)
	time.Sleep(time.Second * 1)
	c2, _ := cpu.Times(false)
	info, _ := disk.IOCounters()
	metrics["cpUser"] = int(c2[0].User - c1[0].User)
	metrics["cpuSys"] = int(c2[0].System - c1[0].System)
	metrics["cpuIowait"] = int(c2[0].Iowait - c1[0].Iowait)
	metrics["iops"] = int(info["disk0"].IopsInProgress)

	NewMysqlMetric := GetMySQLStatus(db)
	for k, v := range NewMysqlMetric {
		metrics[k] = v
	}

	ConditionMap := makeConditionMap(conf)
	result = judge(ConditionMap, metrics)
	return
}

func judge(ConditionMap map[string]int, metrics map[string]int) bool {
	for k, v := range ConditionMap {
		if metrics[k] >= v {
			return true
		}
	}
	return false
}

func makeConditionMap(conf *common.Config) (ConditionMap map[string]int) {
	ConditionMap = make(map[string]int)
	ConditionMap["Cpuser"] = conf.Condition.Cpuser
	ConditionMap["Cpusys"] = conf.Condition.Cpusys
	ConditionMap["Iowait"] = conf.Condition.Iowait
	ConditionMap["Iops"] = conf.Condition.Iops
	ConditionMap["ThreadsRunning"] = conf.Condition.ThreadsRunning
	ConditionMap["ThreadsConnected"] = conf.Condition.ThreadsConnected
	ConditionMap["RowLockWaits"] = conf.Condition.RowLockWaits
	ConditionMap["SlowQuries"] = conf.Condition.SlowQuries
	return
}

func makeSnapshot(childDir string) {
	var lock sync.Mutex
	defer lock.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	lock.Lock()
	go Logio(childDir, &wg)
	go Logmpstat(childDir, &wg)
	wg.Wait()
}
