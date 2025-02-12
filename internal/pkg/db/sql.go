package db

import (
	"database/sql"
	"fmt"
	"pkg/logger"
	"time"

	_ "github.com/lib/pq"
)

type SQLConfig struct {
	User        string `mapstructure:"user"`
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	Port        int16  `mapstructure:"port"`
	DBName      string `mapstructure:"dbName"`
	SSLMode     string `mapstructure:"sslMode"`
	ConnStr     string `mapstructure:"ConnStr"`
	Password    string `mapstructure:"password"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
	MaxLifeTime int    `mapstructure:"maxLifeTime"`
	MaxIdleTime int    `mapstructure:"maxIdleTime"`
}

func (c *SQLConfig) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}

func (c *SQLConfig) GetMySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

func NewConnectPool(dbconf *SQLConfig, log logger.ILogger) *sql.DB {
	db, err := sql.Open("mysql", dbconf.GetMySqlDSN())
	if err != nil {
		log.Errorf("error in connecting to the database %v", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Errorf("unable to connect to the database: %v", err)
		return nil
	}

	log.Infof("database connected successfully")

	db.SetMaxOpenConns(dbconf.MaxOpenConn)
	db.SetMaxIdleConns(dbconf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(dbconf.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(dbconf.MaxIdleTime) * time.Second)
	return db
}
