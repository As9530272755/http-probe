package main

import (
	"flag"
	"http-probe/config"
	"http-probe/http"
	"log"
)

var configFile string

func main() {

	flag.StringVar(&configFile, "c", "http_probe.yaml", "configfile")
	conf, err := config.LoadFile(configFile)
	if err != nil {
		log.Printf("[config.Load.error][err:%v]", err)
		return
	}

	log.Printf("配置是%v", conf)

	go http.StartGin(conf)
	select {}
}
