package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/SultanKs4/ytarchiver/config"
	"github.com/SultanKs4/ytarchiver/domain/datachannel"
	"github.com/SultanKs4/ytarchiver/pkg/fileutils"
	"github.com/SultanKs4/ytarchiver/pkg/ytdata"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func youtubeData(apiKey string, channelId string) (playlist []*youtube.Playlist, playlistItems map[string][]*youtube.PlaylistItem, err error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return
	}
	ytService := ytdata.NewYtData(context.Background(), service)
	playlist, err = ytService.GetChannelPlaylist(channelId)
	if err != nil {
		return
	}
	playlistItems, err = ytService.GetVideosFromPlaylist(playlist)
	return
}

func mappingData(playlist []*youtube.Playlist, playlistItems map[string][]*youtube.PlaylistItem) *datachannel.DataChannel {
	// mapping videos every playlist id
	mapVideoPlaylist := make(map[string][]datachannel.Video)
	for k, v := range playlistItems {
		vids := []datachannel.Video{}
		for _, items := range v {
			thumbnail := items.Snippet.Thumbnails.Maxres
			if thumbnail == nil {
				thumbnail = items.Snippet.Thumbnails.Standard
			}
			vid := datachannel.Video{
				Id:          items.ContentDetails.VideoId,
				PublishedAt: items.ContentDetails.VideoPublishedAt,
				Title:       items.Snippet.Title,
				Thumbnail:   thumbnail.Url}
			vids = append(vids, vid)
		}
		mapVideoPlaylist[k] = append(mapVideoPlaylist[k], vids...)
	}
	// mapping create array playlist
	var pls []datachannel.Playlist
	for _, v := range playlist {
		pl := datachannel.Playlist{
			Title:  v.Snippet.Title,
			Id:     v.Id,
			Videos: mapVideoPlaylist[v.Id]}
		pls = append(pls, pl)
	}
	sort.SliceStable(pls, func(i, j int) bool {
		return len(pls[i].Videos) < len(pls[j].Videos)
	})
	return &datachannel.DataChannel{Playlists: pls, Name: playlist[0].Snippet.ChannelTitle}
}

func saveToJsonFile(dataCh *datachannel.DataChannel, folderPath string) error {
	b, err := json.MarshalIndent(&dataCh, "", " ")
	if err != nil {
		return err
	}

	fp := filepath.Join(fileutils.ResourcePath(), fileutils.SanitizeFilename(folderPath))
	_ = os.MkdirAll(fp, os.ModePerm)
	fileSave := filepath.Join(fp, "channelData.json")

	f, err := os.Create(fileSave)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func Run(config *config.Config, channelId string) error {
	if config.APIKey.Youtube == "" {
		return errors.New("Youtube API Key empty")
	}
	playlist, playlistItems, err := youtubeData(config.APIKey.Youtube, channelId)
	if err != nil {
		return fmt.Errorf("error metadata: %w", err)
	}
	dataCh := mappingData(playlist, playlistItems)
	if err := saveToJsonFile(dataCh, playlist[0].Snippet.ChannelTitle); err != nil {
		return err
	}
	return nil
}
