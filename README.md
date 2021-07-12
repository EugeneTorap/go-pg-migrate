# go-pg-migrate
### It's a very slim golang library for postgres migrations.

## ⚙️ Installation

Make sure you have Go installed ([download](https://golang.org/dl/)). Version `1.16` or higher is required.

```bash
go get github.com/EugeneTorap/go-pg-migrate
```

## Use auto migration in your Go project

* Uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies.
* Uses `io.Reader` streams internally for low memory overhead.
* Thread-safe and no goroutine leaks.
* Auto migration to the latest version at program or server start.

```go
package main

import (
    "embed"
    migrate "github.com/EugeneTorap/go-pg-migrate"
    _ "github.com/EugeneTorap/go-pg-migrate/postgres"
    "github.com/EugeneTorap/go-pg-migrate/source/iofs"
    _ "github.com/lib/pq"
    "log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
    databaseURL := "postgres://postgres:postgres@localhost/db_name?sslmode=disable"
	
    driver, err := iofs.New(embedMigrations, "migrations")
    if err != nil {
		log.Fatal(err)
    }
    migration, err := migrate.NewWithSourceInstance("iofs", driver, databaseURL)
    if err != nil {
		log.Fatal(err)
    }
    _ = migration.Up()
    _, _ = migration.Close()
}
```

## CLI creation a new migration file

```go
package main

import (
	"flag"
	migrate "github.com/EugeneTorap/go-pg-migrate"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		return
	}
	args := flag.Args()[1:]

	switch flag.Arg(0) {
	case "create":
		flagSet := flag.NewFlagSet("create", flag.ExitOnError)
		dirPtr := flagSet.String("dir", "", "Directory to place file in (default: current working directory)")

		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}

		if flagSet.NArg() == 0 {
			log.Fatal("error: please specify name")
		}
		name := flagSet.Arg(0)

		if err := migrate.CreateCmd(*dirPtr, name, true); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Migration file is created")
	os.Exit(1)
}
```

### CLI usage

```bash
go run ./main.go create -dir migrations $(name)
```