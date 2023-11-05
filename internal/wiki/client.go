package wiki

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pthomison/errcheck"
	"github.com/spf13/viper"
	"gitlab.com/tozd/go/mediawiki"
)

func LoadWikiConfig(ctx context.Context) *mediawiki.ProcessDumpConfig {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = 15

	cacheDirPath := viper.GetString("cache.directory")
	cacheFilePath := fmt.Sprintf("%v/%v", cacheDirPath, viper.GetString("cache.filename"))

	err := os.MkdirAll(cacheDirPath, 0755)
	errcheck.Check(err)

	url, err := mediawiki.LatestWikipediaRun(ctx, retryableClient, "enwiki", 0)
	errcheck.Check(err)

	return &mediawiki.ProcessDumpConfig{
		URL:                    url,
		Path:                   cacheFilePath,
		Client:                 retryableClient,
		ItemsProcessingThreads: viper.GetInt("cores"),
		DecompressionThreads:   viper.GetInt("cores"),
		DecodingThreads:        viper.GetInt("cores"),
	}
}
