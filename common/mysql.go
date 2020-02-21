package common

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" //with mysql
)

// NewMySQLConnection the constructor of mysql connecting
func NewMySQLConnection(conf *Config) (mysql.Conn, error) {
	return initMySQLConnection(conf)
}

// QueryResult the result of query
func initMySQLConnection(conf *Config) (db mysql.Conn, err error) {
	db = mysql.New("tcp", "", fmt.Sprintf(
		"%s:%d", conf.DataBase.Host, conf.DataBase.Port),
		conf.DataBase.User, conf.DataBase.Password)
	db.SetTimeout(1000 * time.Millisecond)
	if err = db.Connect(); err != nil {
		err = errors.Wrap(err, "Building mysql connection failed!")
	}
	return
}
