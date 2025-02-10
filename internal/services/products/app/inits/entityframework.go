package inits

import (
	"database/sql"
	"pkg/db"
	"products/ent/gen"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func NewEntClient(sqlConf *db.SQLConfig, sql *sql.DB) *gen.Client {
	drv := entsql.OpenDB(sqlConf.Name, sql)
	return gen.NewClient(gen.Driver(drv))
}
