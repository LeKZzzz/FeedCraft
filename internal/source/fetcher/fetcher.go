package fetcher

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Fetcher handles the I/O, just retrieving the raw binary data.
type Fetcher interface {
	Fetch(ctx context.Context) ([]byte, error)
	BaseURL() string
}

type CachedFetcher struct {
	Internal Fetcher
	Expire   time.Duration
}

func (f *CachedFetcher) BaseURL() string {
	return f.Internal.BaseURL()
}

func (f *CachedFetcher) Fetch(ctx context.Context) ([]byte, error) {
	cacheKey := fmt.Sprintf("%s:%s", constant.PrefixSearchSource, f.BaseURL())

	cached, err := util.CacheGetString(cacheKey)
	if err == nil && cached != "" {
		logrus.WithField("cacheKey", cacheKey).Debugf("search source cache hit")
		return []byte(cached), nil
	}

	start := time.Now()
	data, err := f.Internal.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	logrus.WithField("cacheKey", cacheKey).WithField("duration", time.Since(start)).Infof("search source cache miss, fetched from internal")

	if cacheErr := util.CacheSetString(cacheKey, string(data), f.Expire); cacheErr != nil {
		logrus.WithField("cacheKey", cacheKey).WithError(cacheErr).Warn("failed to cache search source result")
	}
	return data, nil
}
