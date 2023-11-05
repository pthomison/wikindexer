package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pthomison/wikindexer/internal/config"
	"github.com/pthomison/wikindexer/internal/db"
	"github.com/pthomison/wikindexer/internal/types"
	"github.com/pthomison/wikindexer/internal/wiki"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/mediawiki"
	"gitlab.com/tozd/go/x"
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

	BatchDbInsert(ctx, wikiConfig, 5)
}

func BatchDbInsert(ctx context.Context, config *mediawiki.ProcessDumpConfig, batchSize int) {
	articleChan := make(chan mediawiki.Article)

	go func() {
		articles := []db.Insertable{}

		for a := range articleChan {
			articles = append(articles, types.ParseArticle(a))

			if len(articles) > batchSize-1 {
				rows := dbConfig.BatchInsert(articles...)
				incrementorChan <- rows
				articles = []db.Insertable{}
			}
		}

		rows := dbConfig.BatchInsert(articles...)
		incrementorChan <- rows

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

	percentMsg := fmt.Sprintf("Load %3f%% Complete, Approximately %v remaining", percentage, remaingTime)
	rateMsg := fmt.Sprintf("Handling %3f records per second", rate)
	logrus.Info(percentMsg)
	logrus.Info(rateMsg)
}
