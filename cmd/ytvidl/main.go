package main

import (
	"flag"
	"log"

	"github.com/SultanKs4/ytarchiver/config"
	"github.com/SultanKs4/ytarchiver/internal/ytdl"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.LoadConfig()
	if err != nil {
		log.Printf("failed load config.yml: %s, use default value instead", err.Error())
	}

	runJsonFlag := flag.Bool("json", false, "run program by channelData.json inside resource directory")
	idVidFlag := flag.String("id", "", "provide id video for download single video")
	addressVidFlag := flag.String("address", "", "provide address video for download single video")
	flag.Parse()

	if *runJsonFlag {
		err = ytdl.Run(cfg)
	} else {
		if *idVidFlag == "" && *addressVidFlag == "" {
			log.Fatalln("-id and -address cannot empty, must provide one of them")
		}
		err = ytdl.SingleVideo(cfg, *idVidFlag, *addressVidFlag, "")
	}

	if err != nil {
		log.Fatalf(err.Error())
	}
}
