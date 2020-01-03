package common

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

// BaseConf config about dir, log, etc.
type BaseConf struct {
	BaseDir     string
	SnapshotDir string
	SnapshotDay int
	LogDir      string
	LogFile     string
	LogLevel    int
}

// DatabaseConf config about database
type DatabaseConf struct {
	User     string
	Password string
	Host     string
	Port     int
}

// ConditionConf config about conditions
type ConditionConf struct {
	Cpuser           int
	Cpusys           int
	Iowait           int
	Iops             int
	ThreadsRunning   int
	ThreadsConnected int
	RowLockWaits     int
	SlowQuries       int
}

// Config for initializing. This can be loaded from TOML file with -c
type Config struct {
	Base      BaseConf
	DataBase  DatabaseConf
	Condition ConditionConf
}

// NewConfig the constructor of config
func NewConfig(file string) (*Config, error) {
	conf, err := readConf(file)
	return &conf, err
}

func readConf(file string) (conf Config, err error) {
	_, err = os.Stat(file)
	if err != nil {
		file = fmt.Sprint("etc/", file)
		_, err = os.Stat(file)
		if err != nil {
			panic(err)
		}
	}
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}
	snapshotDay, err := cfg.Section("default").Key("snapshot_day").Int()
	if err != nil {
		fmt.Println("No Snapshot!")
		snapshotDay = -1
	}
	logLevel, err := cfg.Section("default").Key("log_level").Int()
	if err != nil {
		fmt.Println("Log level default: 7!")
		logLevel = 7
	}
	host := cfg.Section("mysql").Key("host").String()
	if host == "" {
		fmt.Println("Host default: 127.0.0.1!")
		host = "127.0.0.1"
	}
	snapshotDir := cfg.Section("default").Key("snapshot_dir").String()
	if snapshotDir == "" {
		fmt.Println("SnapshotDir default current dir ")
		snapshotDir = "."
	}
	port, err := cfg.Section("mysql").Key("port").Int()
	if err != nil {
		fmt.Println("Port: default 3306!")
		port = 3306
		err = nil
	}
	cpuser, err := cfg.Section("condtion").Key("cpuser").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		cpuser = 100
		err = nil
	}
	cpusys, err := cfg.Section("condtion").Key("cpusys").Int()
	if err != nil {
		fmt.Println("cpusys: default 100!")
		cpusys = 100
		err = nil
	}
	iowait, err := cfg.Section("condtion").Key("iowait").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		iowait = 100
		err = nil
	}
	iops, err := cfg.Section("condtion").Key("iops").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		iops = 100000
		err = nil
	}
	threadsrunning, err := cfg.Section("condtion").Key("threads_running").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		threadsrunning = 100
		err = nil
	}
	threadsconnected, err := cfg.Section("condtion").Key("threads_running").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		threadsconnected = 10000
		err = nil
	}
	rowlockwaits, err := cfg.Section("condtion").Key("row_lock_waits").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		rowlockwaits = 10000
		err = nil
	}
	slowquries, err := cfg.Section("condtion").Key("slow_quries").Int()
	if err != nil {
		fmt.Println("cpuser: default 100!")
		slowquries = 10000
		err = nil
	}
	conf = Config{
		BaseConf{
			BaseDir:     cfg.Section("default").Key("basedir").String(),
			SnapshotDir: snapshotDir,
			SnapshotDay: snapshotDay,
			LogDir:      cfg.Section("default").Key("log_dir").String(),
			LogFile:     cfg.Section("default").Key("log_file").String(),
			LogLevel:    logLevel,
		},
		DatabaseConf{
			User:     cfg.Section("mysql").Key("user").String(),
			Password: cfg.Section("mysql").Key("password").String(),
			Host:     host,
			Port:     port,
		},
		ConditionConf{
			Cpuser:           cpuser,
			Cpusys:           cpusys,
			Iowait:           iowait,
			Iops:             iops,
			ThreadsRunning:   threadsrunning,
			ThreadsConnected: threadsconnected,
			RowLockWaits:     rowlockwaits,
			SlowQuries:       slowquries,
		},
	}
	return
}
