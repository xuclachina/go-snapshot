package main

import (
	"github.com/ziutek/mymysql/mysql"
)

// GetMySQLStatus for check
func GetMySQLStatus(db mysql.Conn) map[string]int {
	metrics := make(map[string]int)
	rows, _, err := db.Query("SHOW GLOBAL STATUS;")
	if err != nil {
		Log.Alert("get mysql status failed")
	}
	for _, row := range rows {
		switch row.Str(0) {
		case "Threads_running":
			threadRunning, _ := row.Int64Err(1)
			metrics["Threads_running"] = int(threadRunning)
		case "Threads_connected":
			v, _ := row.Int64Err(1)
			metrics["Threads_connected"] = int(v)
		case "Innodb_row_lock_current_waits":
			v, _ := row.Int64Err(1)
			metrics["Innodb_row_lock_current_waits"] = int(v)
		case "Slow_queries":
			v, _ := row.Int64Err(1)
			metrics["Slow_queries"] = int(v)
		}
	}

	row, _, err := db.QueryFirst("SHOW GLOBAL STATUS LIKE 'Slow_queries';")

	if err != nil {
		Log.Alert("get slow queries status failed")
	}
	slowQueries, _ := row.Int64Err(1)
	metrics["Slow_queries"] = int(slowQueries) - metrics["Slow_queries"]

	return metrics
}

// GetInnodbStaus for check
func GetInnodbStaus(db mysql.Conn) (string, error) {
	status, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
	if err != nil {
		Log.Debug("show innodb status error: %+v", err)
		return "", err
	}
	allStatus := status.Str(2)
	return allStatus, nil
}

// GetProcesslist for get mysql processlist
func GetProcesslist(db mysql.Conn) (string, error) {
	processlist, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ FULLL PROCESSLIST")
	if err != nil {
		Log.Debug("show processlist error: %+v", err)
		return "", err
	}
	allProcesslist := processlist.Str(2)
	return allProcesslist, nil
}
