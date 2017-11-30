package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	g "github.com/weihualiu/redapricot/cfg"
	log "github.com/Sirupsen/logrus"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(g.Config().MaxConns)
	DB.SetMaxIdleConns(g.Config().MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
