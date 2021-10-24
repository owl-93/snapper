package cache

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"snapper/model"
	"snapper/utils"
)

var (
	cache = redis.NewClient(&redis.Options{})
)

func IsInitialized() bool {
	return cache != nil
}

/*
	returns a client if initialized otherwise throws an error
*/
func GetInstance() (*redis.Client, error) {
	if !IsInitialized() {
		return nil, errors.New("Cannot get instance of the cache - it is not intializezd")
	}
	return cache, nil
}

/*
	Checks the cache given a URL, if there is an entry it unmarshalls the JSON and returns
	a pointer to the data
*/
func CheckCacheForPage(address string) (*[]model.MetaTag, error) {
	if cache == nil {
		log.Println("error: cache is uninitialized")
		return nil, errors.New("cache not initialized")
	}
	key, e := utils.GetAddressKey(address)
	if e != nil {
		return nil, e
	}
	result, e := cache.Get(cache.Context(), key).Result()
	if e != nil {
		log.Printf("cache miss for %s\n", key)
		return nil, nil
	}
	log.Printf("cache hit for %s\n", key)
	tags := []model.MetaTag{}
	marshallError := json.Unmarshal([]byte(result), &tags)
	if marshallError != nil {
		return nil, e
	}
	return &tags, nil
}

/*
	function to cache the page metadata in the redis cache.
	//TODO: add configurable TTL for cache life
*/
func CachePageMetaData(tags *[]model.MetaTag, address string) error {
	if cache == nil {
		return errors.New("cache unitialized")
	}

	if pageId, err := utils.GetAddressKey(address); err == nil {
		log.Printf("caching meta data for %s\n", pageId)
		if serialized, marshalError := json.Marshal(*tags); marshalError == nil {
			ctx := cache.Context()
			if cacheError := cache.Set(ctx, pageId, string(serialized), time.Hour * 24).Err(); cacheError == nil {
				log.Printf("cached metadata for page %s\n", pageId)
				return nil
			} else {
				return cacheError
			}
		}else {
			log.Println("could not marshal tags into object")
			return marshalError
		}
	} else {
		log.Printf("could not extract page key from address %s\n", address)
		return err
	}
}
