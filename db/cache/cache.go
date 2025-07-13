package cache

import (
	"context"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Feature struct {
	Name          string   `json:"name" redis:"name"`
	IsActive      bool     `json:"is_active" redis:"is_active"`
	BucketsOpened int      `json:"buckets_opened" redis:"buckets_opened"`
	Whitelist     []string `json:"whitelist" redis:"whitelist"`
	Blacklist     []string `json:"blacklist" redis:"blacklist"`
}

func (f *Feature) SetInRedis(ctx context.Context, db *redis.Client) error {
	val := reflect.ValueOf(f).Elem()
	typ := reflect.TypeOf(*f)

	setter := func(p redis.Pipeliner) error {
		featureFields := make([]string, 0)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			redisTag := fieldType.Tag.Get("redis")
			if redisTag == "" {
				continue
			}

			var strVal string
			switch field.Kind() {
			case reflect.String:
				strVal = field.String()
			case reflect.Bool:
				strVal = strconv.FormatBool(field.Bool())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				strVal = strconv.FormatInt(field.Int(), 10)
			case reflect.Slice:
				if field.Type().Elem().Kind() == reflect.String {
					slice := field.Interface().([]string)
					sliceJSON, _ := jsoniter.Marshal(slice)
					strVal = string(sliceJSON)
				} else {
					continue
				}
			default:
				continue
			}

			featureFields = append(featureFields, redisTag)
			featureFields = append(featureFields, strVal)
		}

		if err := db.HSet(ctx, f.Name, featureFields).Err(); err != nil {
			return errors.Wrapf(err, "failed to set key %s", f.Name)
		}
		return nil
	}

	if _, err := db.Pipelined(ctx, setter); err != nil {
		return errors.Wrap(err, "redis pipeline error")
	}
	return nil
}

func (f *Feature) GetByName(ctx context.Context, featureName string, db *redis.Client) error {
	redisData := db.HGetAll(ctx, featureName)
	if err := redisData.Err(); err != nil {
		return errors.Wrapf(err, "failed to get feature by key %s", featureName)
	}

	if len(redisData.Val()) == 0 {
		return nil
	}

	// set all simple types
	if err := redisData.Scan(f); err != nil {
		return errors.Wrapf(err, "failed to scan redis data of feature %s", featureName)
	}

	// set slices
	if whitelistStr, ok := redisData.Val()["whitelist"]; ok {
		_ = jsoniter.Unmarshal([]byte(whitelistStr), &f.Whitelist)
	}
	if blacklistStr, ok := redisData.Val()["blacklist"]; ok {
		_ = jsoniter.Unmarshal([]byte(blacklistStr), &f.Blacklist)
	}

	return nil
}
