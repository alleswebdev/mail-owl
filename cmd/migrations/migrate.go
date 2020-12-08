package main

import (
	"flag"
	"fmt"
	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

func main() {
	flag.Usage = usage
	flag.Parse()

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		Database: os.Getenv("DB_NAME"),
	})
	// check the sql-create command
	checkCreate(db, flag.Args()...)

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func init() {
	currentDir, _ := os.Getwd()
	err := godotenv.Load("/var/external/env")
	if err != nil {
		err = godotenv.Load(path.Join(currentDir, ".env"))
	}
	if err != nil {
		err = godotenv.Load(path.Join(currentDir, "..", "..", ".env"))
	}
	if err != nil {
		err = godotenv.Load()
	}
	if err != nil {
		log.Fatal(err, "Error loading .env file")
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}

func checkCreate(db migrations.DB, a ...string) {
	cmd := "null"

	if len(a) > 0 {
		cmd = a[0]
	}

	if cmd == "create" {

		if len(a) < 2 {
			log.Fatal("please provide migration description")
		}

		version, err := migrations.Version(db)

		filename := fmtMigrationFilename(version+1, strings.Join(a[1:], "_"))
		err = createMigrationFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("created new migration", filename)
		os.Exit(1)
	}

}

func createMigrationFile(filename []string) error {
	basepath, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, name := range filename {
		fname := path.Join(basepath, name)

		_, err = os.Stat(fname)
		if !os.IsNotExist(err) {
			return fmt.Errorf("file=%q already exists", fname)
		}

		err = ioutil.WriteFile(fname, []byte("--put your sql here\n select 1;"), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// Return array of the sql-migration names with current version(up and down)
func fmtMigrationFilename(version int64, descr string) []string {
	//replace all other chars to '_'
	var migrationNameRE = regexp.MustCompile(`[^a-z0-9]+`)

	descr = strings.ToLower(descr)
	descr = migrationNameRE.ReplaceAllString(descr, "_")
	return []string{
		fmt.Sprintf("%d_%s.tx.up.sql", version, descr),
		fmt.Sprintf("%d_%s.tx.down.sql", version, descr)}
}
