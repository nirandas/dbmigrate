## Database migration management for postgres and mysql

    go get -u github.com/nirandas/dbmigrate

**Initialize**

    dbmigrate init

Creates a dbmigrate.json file in the current directory. Edit the file and set the dsn. The type can be either postgres or mysql.

**Generating migrations**

    dbmigrate make hello

Generates a timestamped sql file in the migrations path as configured in dbmigrate.json. You can add the SQL statements to be executed during up and down in this file. Separate multiple statements using --go--

```sql
--dbmigrate:up
create table users (
id bigint unsigned primary key auto_increment,
username varchar(32) not null unique,
email varchar(100) not null unique,
password varchar(100) not null
)engine=innodb charset=utf8;
--dbmigrate:down
drop table users;
```

**Running migrations**


    dbmigrate up

    dbmigrate down

Use -all switch to up or down all migrations.


    dbmigrate up -all

    dbmigrate down -all
