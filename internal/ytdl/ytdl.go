package ytdl

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/SultanKs4/ytarchiver/config"
	"github.com/SultanKs4/ytarchiver/domain/datachannel"
	"github.com/SultanKs4/ytarchiver/pkg/fileutils"
	"github.com/SultanKs4/ytarchiver/pkg/ytvideo"
)

// this will run search channelData.json then download based from every data json
func Run(config *config.Config) error {
	resourcePath := fileutils.ResourcePath()
	fileList, err := fileutils.GetDataJson(resourcePath)
	if err != nil {
		return err
	}

	var channelList []datachannel.DataChannel
	for _, v := range fileList {
		b, err := os.ReadFile(v)
		if err != nil {
			return err
		}
		var dataChannel datachannel.DataChannel
		err = json.Unmarshal(b, &dataChannel)
		if err != nil {
			return err
		}
		channelList = append(channelList, dataChannel)
	}

	for i := range channelList {
		for j := 0; j < len(channelList[i].Playlists); j++ {
			v := channelList[i].Playlists[j]
			for _, vid := range v.Videos {
				vidPath := filepath.Join(resourcePath, fileutils.SanitizeFilename(channelList[i].Name), fileutils.SanitizeFilename(v.Title))
				err := SingleVideo(config, vid.Id, "", vidPath)
				if err != nil {
					if err.Error() == "file exist" {
						log.Print(err)
						continue
					}
					return err
				}
			}
		}
	}

	return nil
}

// Download single video based from id or address
func SingleVideo(config *config.Config, id, address, basePath string) error {
	ytv, err := ytvideo.NewYtVideo(id, address, basePath, config)
	if err != nil {
		return err
	}
	err = ytv.DownloadComposite()
	if err != nil {
		return err
	}
	return nil
}
