package db

import (
	_ "github.com/glebarez/go-sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/pthomison/errcheck"
)

type Client struct {
	DbLocation string

	DB *sqlx.DB
}

type Migratable interface {
	Schema() string
}

func NewClient(DbLocation string) *Client {
	c := &Client{
		DbLocation: DbLocation,
	}

	db, err := sqlx.Open("sqlite", c.DbLocation)
	errcheck.Check(err)

	c.DB = db

	return c
}

func (c *Client) Migrate(migrationStructs ...Migratable) {
	for _, m := range migrationStructs {
		c.DB.MustExec(m.Schema())
	}
}
