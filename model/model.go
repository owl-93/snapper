package model

import (
	"github.com/go-redis/redis/v8"
)

type SnapperConfig struct {
	Port int
	RedisConfig *redis.Options
	DisableCache bool
}

type MetaTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SnapperRequest struct {
	Page string `json:"page"`
	Refresh bool `json:"forceRefresh"`
}
