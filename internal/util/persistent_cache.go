package util

import (
	"FeedCraft/internal/constant"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PersistentCacheEntry is a GORM model for persisting cache entries in SQLite.
type PersistentCacheEntry struct {
	Key   string `gorm:"primaryKey;type:varchar(255)"`
	Value string `gorm:"type:text"`
}

func isPersistentCacheKey(key string) bool {
	return strings.HasPrefix(key, constant.PrefixWebContent)
}

// persistentCacheGet reads a cache entry from SQLite by key.
func persistentCacheGet(key string) (string, bool) {
	db := GetDatabase()
	if db == nil {
		return "", false
	}

	var entry PersistentCacheEntry
	if err := db.Where("key = ?", key).First(&entry).Error; err != nil {
		return "", false
	}
	return entry.Value, true
}

// persistentCacheSet writes a cache entry to SQLite (upsert).
func persistentCacheSet(key, value string) {
	db := GetDatabase()
	if db == nil {
		return
	}

	entry := PersistentCacheEntry{Key: key, Value: value}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&entry).Error
	if err != nil {
		logrus.WithError(err).WithField("key", key).Debug("failed to persist cache entry")
	}
}

// AutoMigratePersistentCache ensures the persistent_cache_entries table exists.
func AutoMigratePersistentCache(db *gorm.DB) {
	if err := db.AutoMigrate(&PersistentCacheEntry{}); err != nil {
		logrus.WithError(err).Warn("failed to migrate persistent_cache_entries table")
	}
}
