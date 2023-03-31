package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/owl-93/snapper/controller"
	"github.com/owl-93/snapper/model"
)

const (
	defaultPort = 8888
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: snapper [-port] <port> [-cache] <redis connection sstring> -no-cache")
	flag.PrintDefaults()
	os.Exit(1)
}

func getAppConfig() *model.SnapperConfig {
	flag.Usage = usage
	config := &model.SnapperConfig{}
	port := flag.Int("port", defaultPort, "port to run snapper on (default 8888)")
	disableCache := flag.Bool("no-cache", false, "disable caching for the application (not recommended)")
	redisUri := flag.String("cache", "", "connection uri for redis cache (defaults to localhost:6379)")
	cacheTtl := flag.Int("cache-ttl", 24, "Number of hours that cache entries are valid for")
	flag.Parse()

	//build config
	config.Port = *port
	redisConfig := &redis.Options{}
	if len(*redisUri) > 0 {
		redisConfig.Addr = *redisUri
	}
	config.RedisConfig = redisConfig
	config.DisableCache = *disableCache
	config.CacheTTL = int64(*cacheTtl)
	log.Printf("configuring snapper with options:\n\nport: %d\ncache disabled: %v\ncache ttl: %d hours\nredis options: %+v\n\n", config.Port, config.DisableCache, config.CacheTTL, *config.RedisConfig)
	return config
}

func main() {

	config := getAppConfig()

	router := gin.Default()

	//attach routers
	controller.InitRoutes(router, config)

	route := fmt.Sprintf("%s:%d", "localhost", config.Port)
	if err := router.Run(route); err != nil {
		panic("unable to start snapper on " + route)
	}

}
