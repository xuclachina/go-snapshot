package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/ziutek/mymysql/mysql"

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

	defer Log.Close()
	// start...
	Log.Info("start snapshot...")
	// go timeout()
	for {
		db, err := common.NewMySQLConnection(conf)
		if err != nil {
			Log.Error("无法建立数据库连接，错误信息：%s", err)
			return
		}
		Log.Info("开始状态检查")
		isnapshot = checkCondition(conf, db)
		now := time.Now().Format("2006-01-02-15-04-05")
		childDir := fmt.Sprintf("%s/%s", conf.Base.SnapshotDir, now)
		if isnapshot {
			Log.Warning("达到触发条件，开始数据库快照信息收集!")
			makeSnapshot(db, childDir)
		}
		db.Close()
		time.Sleep(time.Second * Interval)
	}

}

/*
func timeout() {
	time.AfterFunc(TimeOut*time.Second, func() {
		Log.Error("Execute timeout")
		os.Exit(1)
	})
}*/

func checkCondition(conf *common.Config, db mysql.Conn) (result bool) {
	metrics := make(map[string]int)
	_cpu, _ := cpu.Percent(time.Second, false)
	_info1, _ := disk.IOCounters()

	time.Sleep(time.Second)

	_info2, _ := disk.IOCounters()

	metrics["cpu"] = int(_cpu[0])
	metrics["iops"] = int(_info2[conf.Condition.Device].ReadCount) + int(_info2[conf.Condition.Device].WriteCount) - int(_info1[conf.Condition.Device].ReadCount) - int(_info1[conf.Condition.Device].WriteCount)

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
	ConditionMap["cpu"] = conf.Condition.Cpu
	ConditionMap["Threads_running"] = conf.Condition.ThreadsRunning
	ConditionMap["Threads_connected"] = conf.Condition.ThreadsConnected
	ConditionMap["Innodb_row_lock_current_waits"] = conf.Condition.RowLockWaits
	ConditionMap["Slow_queries"] = conf.Condition.SlowQuries
	return
}

func makeSnapshot(db mysql.Conn, childDir string) {
	err := os.MkdirAll(childDir, 0755)
	if err != nil {
		Log.Alert("创建文件夹失败!")
	}
	var lock sync.Mutex

	//退出释放锁
	defer lock.Unlock()

	var wg sync.WaitGroup
	wg.Add(10)
	lock.Lock()

	//开始记录状态信息
	go LogIo(childDir, &wg)
	go LogMpstat(childDir, &wg)
	go LogDiskSpace(childDir, &wg)
	//TODO:LogMessageInfo
	go LogTop(childDir, &wg)
	//TODO:LogTcpDump
	go LogMemoInfo(childDir, &wg)
	go LogInterrupts(childDir, &wg)
	go LogPs(childDir, &wg)
	go LogNetStat(childDir, &wg)
	go LogInnodbStatus(db, childDir, &wg)
	go LogProcesslist(db, childDir, &wg)
	//TODO:LogTransactions
	//TODO:LogLockInfo
	//TODO:LogSlaveInfo
	//TODO:LogMySQLStatus
	//TODO:LogMySQLVariables

	wg.Wait()
}
