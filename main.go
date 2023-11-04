package main

import (
	"fmt"
	"time"
)

// import (
// 	"context"
// 	"fmt"

// 	_ "github.com/glebarez/go-sqlite"
// 	"github.com/hashicorp/go-retryablehttp"
// 	"github.com/pthomison/errcheck"
// 	"gitlab.com/tozd/go/errors"
// 	"gitlab.com/tozd/go/mediawiki"
// )

// func main() {
// 	ctx := context.TODO()
// 	retryableClient := retryablehttp.NewClient()

// 	// c := db.NewClient(":memory:")

// 	// articles := []wikindexer.Article{}

// 	// c.DB.Select(&articles, "SELECT * FROM place ORDER BY telcode ASC")

// 	url, err := mediawiki.LatestWikipediaRun(ctx, retryableClient, "enwiki", 0)
// 	errcheck.Check(err)

// 	mediawiki.ProcessWikipediaDump(ctx, &mediawiki.ProcessDumpConfig{
// 		URL:    url,
// 		Path:   "./.wikicache",
// 		Client: retryableClient,
// 	}, func(ctx context.Context, a mediawiki.Article) errors.E {
// 		fmt.Printf("%v\n", a.Name)
// 		return nil
// 	})

// }

func main() {
	for {
		fmt.Println("hello world")
		time.Sleep(1 * time.Second)
	}
}
