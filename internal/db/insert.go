package db

import (
	"time"

	sqlite "github.com/glebarez/go-sqlite"
	"github.com/pthomison/errcheck"
	"github.com/sirupsen/logrus"
	sqlite3 "modernc.org/sqlite/lib"
)

type Insertable interface {
	BatchInsert() string
}

func (c *Client) BatchInsert(objs ...Insertable) int64 {
	result, err := c.DB.NamedExec(objs[0].BatchInsert(), objs)
	errcheck.Check(err)

	if sqle, ok := err.(*sqlite.Error); ok {
		if sqle.Code() == sqlite3.SQLITE_BUSY {
			logrus.Info("Sleeping")
			time.Sleep(1 * time.Second)
			c.BatchInsert(objs...)
		} else {
			errcheck.Check(err)
		}
	}

	count, err := result.RowsAffected()
	errcheck.Check(err)

	return count
}
