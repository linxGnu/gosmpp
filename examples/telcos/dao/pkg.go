package dao

import (
	"github.com/linxGnu/gosmpp/examples/telcos/config"

	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var lock sync.RWMutex

// BindDAO init dao
func BindDAO(conf *config.Database) error {
	_db, err := sqlx.Connect(conf.Driver, conf.DSN)
	if err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	db = _db
	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeInMinute) * time.Minute)

	return nil
}

func getDB() *sqlx.DB {
	lock.RLock()
	defer lock.RUnlock()

	return db
}
