package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"snapper/controller"
)

const (
	host = "localhost"
	port = "8888"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: snapper [-p] <port> [-c] <redis connection uri>")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	startingPort := port
	if len(args) == 0 {
		fmt.Printf("no port specified, using default port: %s\n", port)
	} else {
		startingPort = args[0]
	}
	router := gin.Default()

	//attach routers
	controller.InitRoutes(router)

	route := fmt.Sprintf("%s:%s", host, startingPort)
	if err := router.Run(route); err != nil {
		panic("unable to start snapper on " + route)
	}

}
