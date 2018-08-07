package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		if err := client.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGHUP)

	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGKILL:
			os.Exit(0)
		case syscall.SIGHUP:
			cfg, err := readConfig(configPath)
			if err != nil {
				log.Println("cant reload config:", err)
				return
			}
			client.config = cfg
		}
	}
}
