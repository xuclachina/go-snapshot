package main

import (
	"github.com/ziutek/mymysql/mysql"
)

func GetMySQLStatus(db mysql.Conn) map[string]int {
	metrics := make(map[string]int)
	rows, _, err := db.Query("SHOW GLOBAL STATUS;")
	if err != nil {
		Log.Alert("get mysql status failed")
	}
	for _, row := range rows {
		switch row.Str(0) {
		case "Threads_running":
			threadRunning, _ := row.Int64Err(0)
			metrics["Threads_running"] = int(threadRunning)
		case "Threads_connected":
			v, _ := row.Int64Err(0)
			metrics["Threads_connected"] = int(v)
		case "Innodb_row_lock_current_waits":
			v, _ := row.Int64Err(0)
			metrics["Innodb_row_lock_current_waits"] = int(v)
		case "Slow_queries":
			v, _ := row.Int64Err(0)
			metrics["Slow_queries"] = int(v)
		}
	}

	row, _, err := db.QueryFirst("SHOW GLOBAL STATUS LIKE 'Slow_queries';")

	if err != nil{
		Log.Alert("get slow queries status failed")
	}
	slowQueries, _ := row.Int64Err(0)
	metrics["Slow_queries"] = int(slowQueries) - metrics["Slow_queries"]

	return metrics
}
