package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/anyandrea/smartdev/lib/env"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

type Adapter interface {
	GetDatabase() *sql.DB
	GetURI() string
	GetType() string
	RunMigrations(string) error
}

func NewAdapter() (db Adapter) {
	var databaseType, databaseUri string

	// get db type
	databaseType = env.Get("WEATHERDB_TYPE", "postgres")

	// check for VCAP_SERVICES first
	vcap, err := cfenv.Current()
	if err != nil {
		log.Println("Could not parse VCAP environment variables")
		log.Println(err)
	} else {
		service, err := vcap.Services.WithName("smartdevdb")
		if err != nil {
			log.Println("Could not find weatherdb service in VCAP_SERVICES")
			log.Fatal(err)
		}
		databaseUri = fmt.Sprintf("%v", service.Credentials["uri"])

		// stupid servicebroker is giving us an improperly formatted DSN
		if databaseType == "mysql" {
			databaseUri = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?multiStatements=true&parseTime=true",
				service.Credentials["username"],
				service.Credentials["password"],
				service.Credentials["hostname"],
				service.Credentials["port"],
				service.Credentials["database"])
		}
	}

	// if WEATHERDB_URI is not yet set then try to read it from ENV
	if len(databaseUri) == 0 {
		databaseUri = env.MustGet("WEATHERDB_URI")
	}

	// setup database adapter
	switch databaseType {
	case "postgres":
		db = newPostgresAdapter(databaseUri)
	case "mysql":
		db = newMysqlAdapter(databaseUri)
	case "sqlite":
		db = newSQLiteAdapter(databaseUri)
	default:
		log.Fatalf("Invalid database type: %s\n", databaseType)
	}

	// panic if no database adapter was set up
	if db == nil {
		log.Fatal("Could not set up database adapter")
	}

	return db
}
