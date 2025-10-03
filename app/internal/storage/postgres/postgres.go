package postgres

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

const (
	PicturesTable = "pictures"
)

var (
	strategy = retry.Strategy{
		Attempts: 5,
		Delay:    1,
		Backoff:  1.25,
	}
)

type Postgres struct {
	db *dbpg.DB
}

func New(conn string) *Postgres {
	db, err := dbpg.New(conn, nil, &dbpg.Options{
		MaxOpenConns:    100,
		MaxIdleConns:    20,
		ConnMaxLifetime: 1 * time.Hour,
	})
	if err != nil {
		panic(err)
	}

	err = db.Master.PingContext(context.Background())
	if err != nil {
		panic(err)
	}

	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Shutdown() {
	_ = p.db.Master.Close()
}
