package dbmigrate

import (
	"database/sql"
)

type MigrationRecord struct {
	Migration string
	Batch     int
}

func (mc *MigrationContent) RunUp(db *sql.DB, batch int) error {
	println(mc.Up)
	_, err := db.Exec(mc.Up)
	if err != nil {
		return err
	}
sql:="insert into _migration_history_ (migration,batch) values($1,$2)"
if Config.DBType == "mysql"{
sql = "insert into _migration_history_ (migration,batch) values(?,?)"
}
	_, err = db.Exec(sql, mc.Name, batch)
	return nil
}

func (mc *MigrationContent) RunDown(db *sql.DB) error {
	_, err := db.Exec(mc.Down)
	if err != nil {
		return err
	}
sql:="delete from _migration_history_ where migration = $1"
if Config.DBType == "mysql"{
sql = "delete from _migration_history_ where migration = ?"
}
	db.Exec(sql, mc.Name)
	return nil
}
