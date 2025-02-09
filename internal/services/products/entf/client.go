package entf

import (
	"database/sql"
	"pkg/db"
	"products/entf/generated"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func NewEntClient(sqlConf *db.SQLConfig, sql *sql.DB) *generated.Client {
	drv := entsql.OpenDB(sqlConf.Name, sql)
	return generated.NewClient(generated.Driver(drv))
}
