package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/peedans/beerleo/config"
	"github.com/peedans/beerleo/modules/servers"
	"github.com/peedans/beerleo/pkg/databases"
	"os"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.dev"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	servers.NewServer(cfg, db).Start()

}
