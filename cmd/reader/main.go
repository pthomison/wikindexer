package main

import (
	"context"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pthomison/errcheck"
	"github.com/pthomison/wikindexer/internal/config"
	"github.com/sirupsen/logrus"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/mediawiki"
	"gitlab.com/tozd/go/x"
)

func init() {
	config.ViperInit()
}

func main() {
	ctx := context.TODO()
	retryableClient := retryablehttp.NewClient()

	// cacheDirPath := viper.GetString("cache.directory")
	// cacheFilePath := fmt.Sprintf("%v/%v", cacheDirPath, viper.GetString("cache.filename"))

	// err := os.MkdirAll(cacheDirPath, 0755)
	// errcheck.Check(err)

	url, err := mediawiki.LatestWikipediaRun(ctx, retryableClient, "enwiki", 0)
	errcheck.Check(err)

	mediawiki.ProcessWikipediaDump(ctx, &mediawiki.ProcessDumpConfig{
		URL:      url,
		Path:     "./tmp",
		Client:   retryableClient,
		Progress: LogProgress,
	}, func(ctx context.Context, a mediawiki.Article) errors.E {
		return nil
	})
}

func LogProgress(ctx context.Context, p x.Progress) {
	logrus.Info(p)
}
