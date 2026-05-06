package craft

import (
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"encoding/json"

	"github.com/gorilla/feeds"
)

// CacheKeyMeta holds the metadata used to generate a cache key.
type CacheKeyMeta struct {
	Title string `json:"title"`
	Id    string `json:"id"`
	Link  string `json:"link"`
}

// CacheEntry is the structured JSON value stored in Redis.
type CacheEntry struct {
	Value string       `json:"value"`
	Meta  CacheKeyMeta `json:"meta"`
}

// CacheKeyResult is the return value of a cache key generator.
type CacheKeyResult struct {
	Hash     string
	MetaJSON string
}

// unifiedItemKeyGen generates a cache key from a feeds.Item using MD5(Title + Id + Link).
func unifiedItemKeyGen(item *feeds.Item) (CacheKeyResult, error) {
	link := ""
	if item.Link != nil {
		link = item.Link.Href
	} else if item.Source != nil {
		link = item.Source.Href
	}
	meta := CacheKeyMeta{Title: item.Title, Id: item.Id, Link: link}
	metaJSON, _ := json.Marshal(meta)
	return CacheKeyResult{
		Hash:     util.GetMD5Hash(item.Title + item.Id + link),
		MetaJSON: string(metaJSON),
	}, nil
}

// unifiedArticleKeyGen generates a cache key from a CraftArticle using MD5(Title + Id + Link).
func unifiedArticleKeyGen(article *model.CraftArticle) (CacheKeyResult, error) {
	meta := CacheKeyMeta{Title: article.Title, Id: article.Id, Link: article.Link}
	metaJSON, _ := json.Marshal(meta)
	return CacheKeyResult{
		Hash:     util.GetMD5Hash(article.Title + article.Id + article.Link),
		MetaJSON: string(metaJSON),
	}, nil
}
