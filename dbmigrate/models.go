package dbmigrate

import (
	"database/sql"
)

type MigrationRecord struct {
	Migration string
	Batch     int
}

func (mc *MigrationContent) RunUp(db *sql.DB, batch int) error {
	for _, up := range mc.Up {
		println(up)
		_, err := db.Exec(up)
		if err != nil {
			return err
		}
	}
	sql := "insert into _migration_history_ (migration,batch) values($1,$2)"
	if Config.DBType == "mysql" {
		sql = "insert into _migration_history_ (migration,batch) values(?,?)"
	}
	_, _ = db.Exec(sql, mc.Name, batch)
	return nil
}

func (mc *MigrationContent) RunDown(db *sql.DB) error {
	for _, down := range mc.Down {
		_, err := db.Exec(down)
		if err != nil {
			return err
		}
	}
	sql := "delete from _migration_history_ where migration = $1"
	if Config.DBType == "mysql" {
		sql = "delete from _migration_history_ where migration = ?"
	}
	db.Exec(sql, mc.Name)
	return nil
}
