package db

import (
	"fmt"
	"strings"

	_ "github.com/glebarez/go-sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/pthomison/errcheck"
)

type Client struct {
	DbLocation string

	DB *sqlx.DB
}

func NewClient(DbLocation string) *Client {
	c := &Client{
		DbLocation: DbLocation,
	}

	sqliteOptions := map[string]string{
		"journal_mode": "WAL",
		"synchronous":  "normal",
	}
	pragmaArr := []string{}

	for option, val := range sqliteOptions {
		pragmaArr = append(pragmaArr, fmt.Sprintf("_pragma=%v(%v)", option, val))
	}

	pragmaStr := strings.Join(pragmaArr, "&")

	db, err := sqlx.Open("sqlite", fmt.Sprintf("%v?%v", c.DbLocation, pragmaStr))
	errcheck.Check(err)

	c.DB = db

	return c
}
