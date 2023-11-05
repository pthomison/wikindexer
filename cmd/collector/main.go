package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pthomison/errcheck"
	"github.com/pthomison/wikindexer/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/mediawiki"
	"gitlab.com/tozd/go/x"
)

func init() {
	config.ViperInit()
	logrus.Info(viper.AllSettings())
	// viper.Set("cache.directory", ".")
	// viper.Set("cache.filename", "temp")
}

func main() {
	logrus.Info("Starting WikIndexer Collector")

	ctx := context.TODO()
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = 15

	cacheDirPath := viper.GetString("cache.directory")
	cacheFilePath := fmt.Sprintf("%v/%v", cacheDirPath, viper.GetString("cache.filename"))

	err := os.MkdirAll(cacheDirPath, 0755)
	errcheck.Check(err)

	url, err := mediawiki.LatestWikipediaRun(ctx, retryableClient, "enwiki", 0)
	errcheck.Check(err)

	mediawiki.ProcessWikipediaDump(ctx, &mediawiki.ProcessDumpConfig{
		URL:                    url,
		Path:                   cacheFilePath,
		Client:                 retryableClient,
		Progress:               LogProgress,
		ItemsProcessingThreads: viper.GetInt("cores"),
		DecompressionThreads:   viper.GetInt("cores"),
		DecodingThreads:        viper.GetInt("cores"),
	}, func(ctx context.Context, a mediawiki.Article) errors.E {
		return nil
	})
}

func LogProgress(ctx context.Context, p x.Progress) {
	percentage := p.Percent()
	remaingTime := time.Until(p.Estimated()).Truncate(time.Second).String()

	logMsg := fmt.Sprintf("Download %3f%% Complete, Approximately %v remaining", percentage, remaingTime)
	logrus.Info(logMsg)
}
