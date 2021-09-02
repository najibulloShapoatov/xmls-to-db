package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//"github.com/jinzhu/gorm"
	//github.com/jinzhu/gorm/dialects/postgres .
	//_ "github.com/jinzhu/gorm/dialects/postgres"

	//github.com/lib/pq
	_ "github.com/lib/pq"
)

//Config database struct
type Config struct {
	Host            string
	Port            string
	Dbname          string
	SslMode         string
	User            string
	Pass            string
	ConnMaxLifetime int
	MaxOpenConns    int
	MaxIdleConns    int
	ApplicationName string
}

var db *sql.DB
var err error

//Init - Database init
func Init(dbConf *Config) {

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable  application_name=%s",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Port,
		dbConf.Dbname,
		dbConf.ApplicationName,
	)
	//db, err = gorm.Open("postgres", dbinfo)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(dbConf.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(dbConf.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbConf.ConnMaxLifetime))
}

//GetDB - get DB
func GetDB() *sql.DB {
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
