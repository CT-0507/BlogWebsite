package cache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	CacheHit  atomic.Uint64
	CacheMiss atomic.Uint64
)

func NewRedisClient() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis")

	return rdb
}

// func BuildCacheKey(prefix string, params map[string]interface{}) string {
// 	// Sort keys to ensure deterministic cache keys
// 	keys := make([]string, 0, len(params))

// 	for k := range params {
// 		keys = append(keys, k)
// 	}

// 	sort.Strings(keys)

// 	parts := make([]string, 0, len(keys))

// 	for _, k := range keys {
// 		v := params[k]

// 		// skip nil values
// 		if v == nil {
// 			continue
// 		}

// 		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
// 	}

// 	raw := strings.Join(parts, ":")

// 	return fmt.Sprintf("%s:%s", prefix, raw)
// }

func normalizeValue(v any) any {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		return rv.Elem().Interface()
	}

	return v
}

func BuildCacheKey(
	prefix string,
	params map[string]any,
	readableParams map[string]any,
) string {
	// ---------- Build deterministic hash ----------
	paramKeys := make([]string, 0, len(params))

	for k := range params {
		paramKeys = append(paramKeys, k)
	}

	sort.Strings(paramKeys)

	var hashBuilder strings.Builder

	for _, k := range paramKeys {
		v := params[k]

		if v == nil {
			continue
		}

		fmt.Fprintf(&hashBuilder, "%s=%v:", k, normalizeValue(v))
	}

	hash := sha1.Sum([]byte(hashBuilder.String()))
	hashString := hex.EncodeToString(hash[:8])

	// ---------- Build readable section ----------
	readableKeys := make([]string, 0, len(readableParams))

	for k := range readableParams {
		readableKeys = append(readableKeys, k)
	}

	sort.Strings(readableKeys)

	readableParts := make([]string, 0, len(readableKeys))

	for _, k := range readableKeys {
		v := readableParams[k]

		if v == nil {
			continue
		}

		readableParts = append(
			readableParts,
			fmt.Sprintf("%s=%v", k, normalizeValue(v)),
		)
	}

	// ---------- Final key ----------
	if len(readableParts) > 0 {
		return fmt.Sprintf(
			"%s:%s:%s",
			prefix,
			strings.Join(readableParts, ":"),
			hashString,
		)
	}

	return fmt.Sprintf(
		"%s:%s",
		prefix,
		hashString,
	)
}

func GetRedisHitRate(ctx context.Context, rdb *redis.Client) (float64, error) {
	info, err := rdb.Info(ctx, "stats").Result()
	if err != nil {
		return 0, err
	}

	var hits, misses float64

	lines := strings.SplitSeq(info, "\n")

	for line := range lines {
		if strings.HasPrefix(line, "keyspace_hits:") {
			fmt.Sscanf(line, "keyspace_hits:%f", &hits)
		}

		if strings.HasPrefix(line, "keyspace_misses:") {
			fmt.Sscanf(line, "keyspace_misses:%f", &misses)
		}
	}

	total := hits + misses

	if total == 0 {
		return 0, nil
	}

	return (hits / total) * 100, nil
}
