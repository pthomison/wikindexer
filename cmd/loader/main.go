package main

import (
	"context"
	"fmt"
	"time"

	"github.com/glebarez/go-sqlite"
	"github.com/pthomison/errcheck"
	"github.com/pthomison/wikindexer/internal/config"
	"github.com/pthomison/wikindexer/internal/db"
	"github.com/pthomison/wikindexer/internal/types"
	"github.com/pthomison/wikindexer/internal/wiki"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/mediawiki"
	"gitlab.com/tozd/go/x"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	incrementorChan = make(chan int64)
	processChan     = make(chan x.Progress)
	dbConfig        *db.Client
)

func init() {
	config.ViperInit()
	logrus.Info(viper.AllSettings())

	dbConfig = db.NewClient(fmt.Sprintf("%v/%v", viper.GetString("cache.directory"), "sqlite.db"))
	dbConfig.Migrate(&types.Article{})
}

func main() {
	logrus.Info("Starting WikIndexer Reader")

	ctx := context.TODO()

	wikiConfig := wiki.LoadWikiConfig(ctx)
	wikiConfig.Progress = ProgressPassthrough

	go ProgressCounter()

	BatchInsertExecution(ctx, wikiConfig, 5)
}

func AttemptInsert(a mediawiki.Article) {
	res, err := dbConfig.DB.Exec(`INSERT OR IGNORE INTO article (id, name, url, abstract)
		VALUES ($1, $2, $3, $4)`, a.Identifier, a.Name, a.URL, a.Abstract)

	if sqle, ok := err.(*sqlite.Error); ok {
		if sqle.Code() == sqlite3.SQLITE_BUSY {
			logrus.Info("Sleeping")
			time.Sleep(1 * time.Second)
			AttemptInsert(a)
		} else {
			errcheck.Check(err)
		}
	}

	count, err := res.RowsAffected()
	errcheck.Check(err)

	incrementorChan <- count
}

func AttemptBatch(articles []mediawiki.Article) {
	result, err := dbConfig.DB.NamedExec(`INSERT OR IGNORE INTO article (id, name, url, abstract)
	VALUES (:identifier, :name, :url, :abstract)`, articles)
	errcheck.Check(err)

	if sqle, ok := err.(*sqlite.Error); ok {
		if sqle.Code() == sqlite3.SQLITE_BUSY {
			logrus.Info("Sleeping")
			time.Sleep(1 * time.Second)
			AttemptBatch(articles)
		} else {
			errcheck.Check(err)
		}
	}

	count, err := result.RowsAffected()
	errcheck.Check(err)

	incrementorChan <- count
}

func SingleInsertExecution(ctx context.Context, config *mediawiki.ProcessDumpConfig) {
	mediawiki.ProcessWikipediaDump(ctx, config, func(ctx context.Context, a mediawiki.Article) errors.E {

		AttemptInsert(a)

		return nil
	})
}

func BatchInsertExecution(ctx context.Context, config *mediawiki.ProcessDumpConfig, batchSize int) {
	articleChan := make(chan mediawiki.Article)

	go func() {
		articles := []mediawiki.Article{}

		for a := range articleChan {
			articles = append(articles, a)

			if len(articles) > batchSize-1 {
				AttemptBatch(articles)

				articles = []mediawiki.Article{}
			}
		}
	}()

	mediawiki.ProcessWikipediaDump(ctx, config, func(ctx context.Context, a mediawiki.Article) errors.E {
		articleChan <- a
		return nil
	})
}

func ProgressCounter() {
	var total int64 = 0

	for {
		select {
		case i := <-incrementorChan:
			total += i

		case p := <-processChan:
			LogProgress(p, total)
		}
	}
}

func ProgressPassthrough(ctx context.Context, p x.Progress) {
	processChan <- p
}

func LogProgress(p x.Progress, processedRecords int64) {
	percentage := p.Percent()
	remaingTime := time.Until(p.Estimated()).Truncate(time.Second).String()

	rate := float64(processedRecords) / p.Elapsed.Seconds()

	logrus.Info("Records: ", processedRecords)
	logrus.Info("Seconds: ", p.Elapsed.Seconds())

	logMsg := fmt.Sprintf("Load %3f%% Complete, Approximately %v remaining", percentage, remaingTime)
	logrus.Info(logMsg)
	logrus.Info(fmt.Sprintf("Handling %3f records per second", rate))
}
