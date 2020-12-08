### Миграции

Для миграций используется пакет [go-pg migrations](github.com/go-pg/migrations/v7)

Флаги запуска migrations:
```bash
 - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
```

Пример:
```bash
> go run *.go init
version is 0

> go run *.go version
version is 0

> go run *.go
creating table my_table...
adding id column...
seeding my_table...
migrated from version 0 to 4

> go run *.go version
version is 4

> go run *.go reset
truncating my_table...
dropping id column...
dropping table my_table...
migrated from version 4 to 0

> go run *.go up 2
creating table my_table...
adding id column...
migrated from version 0 to 2

> go run *.go
seeding my_table...
migrated from version 2 to 4

> go run *.go down
truncating my_table...
migrated from version 4 to 3

> go run *.go version
version is 3

> go run *.go set_version 1
migrated from version 3 to 1

> go run *.go create add email to users
created new migration [2_add_email_to_users.tx.up.sql 2_add_email_to_users.tx.down.sql]
```

Для файлов миграций используется тип sql-migrations, миграции пишутся на голом SQL:
```
SQL migrations are automatically picked up if placed in the same folder with main.go or Go migrations. SQL migrations must have one of the following extensions:

    .up.sql - up migration;
    .down.sql - down migration;
    .tx.up.sql - transactional up migration;
    .tx.down.sql - transactional down migration.

```

Создать миграцию:
```bash
go run migrate.go create [имя миграции]
```

Пример:

```bash
alles@ubuntu:~/go/producer/cmd/migrations$ go run migrate.go create create sms table
created new migration [2_create_sms_table.tx.up.sql 2_create_sms_table.tx.down.sql]
exit status 1
```