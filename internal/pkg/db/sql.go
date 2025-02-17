package db

import (
	"context"
	"database/sql"
	"fmt"
	"pkg/logger"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

func NewConnectPool(ctx context.Context, dbconf *SQLConfig, log logger.Zapper) *sql.DB {
	db, err := sql.Open("mysql", dbconf.GetMySqlDSN())
	if err != nil {
		log.Errorf(ctx, "error in connecting to the database %v", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Errorf(ctx, "unable to connect to the database: %v", err)
		return nil
	}

	log.Infof(ctx, "database connected successfully")

	db.SetMaxOpenConns(dbconf.MaxOpenConn)
	db.SetMaxIdleConns(dbconf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(dbconf.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(dbconf.MaxIdleTime) * time.Second)
	return db
}

func NewOtelDBConnectionPool(ctx context.Context, dbconf *SQLConfig, log logger.Zapper) *sql.DB {

	driverName, err := otelsql.Register("mysql", otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))

	if err != nil {
		log.Errorf(ctx, "error in registering OpenTelemetry SQL driver: %v", err)
		return nil
	}

	// Open the database connection using the wrapped driver
	db, err := sql.Open(driverName, dbconf.GetMySqlDSN())
	if err != nil {
		log.Errorf(ctx, "error in connecting to the database %v", err)
		return nil
	}

	// Register DB stats metrics
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))

	if err != nil {
		log.Errorf(ctx, "error in registering DB stats metrics %v", err)
		return nil
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Errorf(ctx, "unable to connect to the database: %v", err)
		return nil
	}

	log.Info(ctx, "database connected successfully %s")

	// Set database connection pool parameters
	db.SetMaxOpenConns(dbconf.MaxOpenConn)
	db.SetMaxIdleConns(dbconf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(dbconf.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(dbconf.MaxIdleTime) * time.Second)
	return db
}
