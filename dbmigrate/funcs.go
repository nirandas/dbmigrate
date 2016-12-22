package dbmigrate

import (
	"database/sql"
	_ "github.com/lib/pq"
_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func ListMigrations(migration string) ([]string, []string) {
	prev := []string{}
	next := []string{}
	fs, err := ioutil.ReadDir(Config.Path)
	if err != nil {
		log.Fatalf("Failed to parse migrations %s", err.Error())
	}
	found := false
	for _, f := range fs {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		p := path.Join(Config.Path, f.Name())
		if found {
			next = append(next, p)
		} else {
			prev = append(prev, p)
		}
		if f.Name() == migration {
			found = true
		}
	}

	if !found && len(next) == 0 && len(prev) > 0 {
		return next, prev
	}
	return prev, next
}

func LastMigration(rec []MigrationRecord) MigrationRecord {
	lm := MigrationRecord{}
	if len(rec) > 0 {
		lm = rec[len(rec)-1]
	}
	return lm
}

func OpenDb() *sql.DB {
	if Config.DBType == "" {
		log.Fatalf("Database type not set in configuration, set the 'type' key")
	}
	if Config.DSN == "" {
		log.Fatalf("Database dsn not set in configuration, set the 'dsn' key")
	}

	db, err := sql.Open(Config.DBType, Config.DSN)
	if err != nil {
		log.Fatalf("Couldn't open the database %s", err.Error())
	}
	db.Exec("create table _migration_history_ (migration varchar (255), batch int);")
	return db
}

func ParseMigrationRecords(db *sql.DB) []MigrationRecord {
	rec := []MigrationRecord{}
	rows, err := db.Query("select * from _migration_history_ order by migration")
	if err != nil {
		log.Fatalf("Reading migration history failed %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		r := MigrationRecord{}
		err = rows.Scan(&r.Migration, &r.Batch)
		if err != nil {
			log.Fatalf("Reading migration history failed %s", err.Error())
		}
		rec = append(rec, r)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Reading migration history failed %s", err.Error())
	}
	return rec
}

type MigrationContent struct {
	Name string
	Up   string
	Down string
}

func ParseMigration(mpath string) *MigrationContent {
	data, err := ioutil.ReadFile(mpath)
	if err != nil {
		log.Fatalf("Failed parsing migration %s", mpath, err.Error())
	}
	content := string(data)
	upIdx := strings.Index(content, "--dbmigrate:up")
	downIdx := strings.Index(content, "--dbmigrate:down")
	if upIdx == -1 || downIdx == -1 {
		log.Fatalf("Invalid migration file %s", mpath)
	}

	content = strings.Replace(content, "--dbmigrate:up", "", 1)
	parts := strings.Split(content, "--dbmigrate:down")
	if len(parts) != 2 {
		log.Fatalf("Invalid migration file %s", mpath)
	}
	return &MigrationContent{
		Up:   parts[0],
		Down: parts[1],
		Name: path.Base(mpath),
	}
}
