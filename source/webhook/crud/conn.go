package crud

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
	"webhook/config"
	"webhook/log"
)

var (
	db          *sqlx.DB
	pingTimeout = time.Second
	pingRate    = time.Second * 10
)

func init() {

	var err error

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",

		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Addr,
		config.Config.Database.DatabaseName,
	)

	db = sqlx.MustOpen("postgres", dbURL)

	err = ping()

	if err != nil {
		log.Error(fmt.Sprintf("unable to ping database: %s", err.Error()))
		panic(err)
	}

	go func() {

		for {

			time.Sleep(pingRate)

			err = ping()
			if err != nil {

				log.Error(fmt.Sprintf("unable to ping database: %s", err.Error()))

			}

		}

	}() //checks database availability

}

func ping() (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	err = db.PingContext(ctx)

	return

}
