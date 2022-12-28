package main

import (
	"flag"
	"log"

	"github.com/SultanKs4/ytarchiver/config"
	"github.com/SultanKs4/ytarchiver/internal/metadata"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	channelIdFlag := flag.String("chid", "", "youtube channel id")
	flag.Parse()
	if *channelIdFlag == "" {
		log.Fatalf("error: chid flag empty")
	}

	err = metadata.Run(cfg, *channelIdFlag)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
