package redis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

func LoadSegments(path string, variant int, isAd bool) error {
	var key string
	if isAd {
		key = "ad_segments"
	} else {
		key = fmt.Sprintf("video_segments_%d", variant)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".ts") {
			rdb.RPush(ctx, key, entry.Name())
		}
	}
	return nil
}

func PopSegments(variant int, count int) ([]string, error) {
	key := fmt.Sprintf("video_segments_%d", variant)
	vals, err := rdb.LRange(ctx, key, 0, int64(count-1)).Result()
	if err != nil {
		return nil, err
	}
	rdb.LTrim(ctx, key, int64(count), -1)
	return vals, nil
}

func GetAdSegment() (string, error) {
	vals, err := rdb.LRange(ctx, "ad_segments", 0, 0).Result()
	if err != nil || len(vals) == 0 {
		return "", err
	}
	return vals[0], nil
}

func ShouldInsertAd() bool {
	last, err := rdb.Get(ctx, "last_switch").Int64()
	if err != nil || time.Now().Unix()-last >= 30 {
		rdb.Set(ctx, "last_switch", time.Now().Unix(), 0)
		return true
	}
	return false
}

func FlushAll() {
	_ = rdb.FlushAll(ctx).Err()
}

func GetCurrentMode() string {
	mode, err := rdb.Get(ctx, "mode").Result()
	if err != nil {
		return "content"
	}
	return mode
}

func GetVariantIndex(variant int) int {
	key := fmt.Sprintf("index_%d", variant)
	val, err := rdb.Get(ctx, key).Int()
	if err != nil {
		return 0
	}
	return val
}

func IncrementVariantIndex(variant int, by int) {
	key := fmt.Sprintf("index_%d", variant)
	_ = rdb.IncrBy(ctx, key, int64(by)).Err()
}

func GetSegmentsForVariant(variant, start, count int) ([]string, error) {
	key := fmt.Sprintf("video_segments_%d", variant)
	return rdb.LRange(ctx, key, int64(start), int64(start+count-1)).Result()
}
