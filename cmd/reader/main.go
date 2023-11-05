package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pthomison/wikindexer/internal/config"
	"github.com/pthomison/wikindexer/internal/wiki"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/mediawiki"
	"gitlab.com/tozd/go/x"
)

func init() {
	config.ViperInit()
	logrus.Info(viper.AllSettings())
}

func main() {
	logrus.Info("Starting WikIndexer Reader")

	ctx := context.TODO()
	wikiConfig := wiki.LoadWikiConfig(ctx)
	wikiConfig.Progress = LogProgress

	mediawiki.ProcessWikipediaDump(ctx, wikiConfig, func(ctx context.Context, a mediawiki.Article) errors.E {
		// logrus.Info(a.Name)
		return nil
	})
}

func LogProgress(ctx context.Context, p x.Progress) {
	percentage := p.Percent()
	remaingTime := time.Until(p.Estimated()).Truncate(time.Second).String()

	logMsg := fmt.Sprintf("Read %3f%% Complete, Approximately %v remaining", percentage, remaingTime)
	logrus.Info(logMsg)
}
