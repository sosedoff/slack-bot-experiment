package main

import (
	"flag"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./config.yml", "Path to config file")
	flag.Parse()
}

func main() {
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(config)
	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
