package ytdata

import (
	"context"
	"log"

	"google.golang.org/api/youtube/v3"
)

type YtData struct {
	ctx      context.Context
	service  *youtube.Service
	basePart []string
}

func NewYtData(ctx context.Context, service *youtube.Service) *YtData {
	return &YtData{ctx: ctx, service: service, basePart: []string{"snippet", "contentDetails"}}
}

// TODO: WIP
func (ytd *YtData) GetChannel(username string) (*youtube.ChannelListResponse, error) {
	call := ytd.service.Channels.List(ytd.basePart).ForUsername(username)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ytd *YtData) GetChannelPlaylist(channelId string) ([]*youtube.Playlist, error) {
	call := ytd.service.Playlists.List(ytd.basePart).ChannelId(channelId).MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (ytd *YtData) GetVideosFromPlaylist(playlists []*youtube.Playlist) (map[string][]*youtube.PlaylistItem, error) {
	playlistItems := make(map[string][]*youtube.PlaylistItem)
	for _, vPlaylist := range playlists {
		log.Printf("get items from playlist: %s", vPlaylist.Snippet.Title)
		call := ytd.service.PlaylistItems.List(ytd.basePart).PlaylistId(vPlaylist.Id).MaxResults(50)
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}
		items := resp.Items
		for resp.NextPageToken != "" {
			call = call.PageToken(resp.NextPageToken)
			resp, err = call.Do()
			if err != nil {
				return nil, err
			}
			items = append(items, resp.Items...)
		}
		for _, vItem := range items {
			switch vItem.Snippet.Title {
			case "Deleted video", "Private video":
				continue
			default:
				playlistItems[vPlaylist.Id] = append(playlistItems[vPlaylist.Id], vItem)
			}
		}
	}
	return playlistItems, nil
}
