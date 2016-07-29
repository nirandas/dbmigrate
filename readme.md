##Database migration for postgres

    go get github.com/nirandas/dbmigrate

**Initialize**

    dbmigrate init

Creates a dbmigrate.json file in the current directory. Edit the file and set the dsn.

**Generating migrations**

    dbmigrate make hello

Generates a timestamped file in the migrations path as configured in dbmigrate.json

**Running migrations**

    dbmigrate up

    dbmigrate down

Use -all switch to up or down all migrations.

    dbmigrate up -all

    dbmigrate down -all
