package datachannel

type DataChannel struct {
	Name      string     `json:"name"`
	Playlists []Playlist `json:"playlist"`
}

type Playlist struct {
	Title  string  `json:"title"`
	Id     string  `json:"id"`
	Videos []Video `json:"videos"`
}

type Video struct {
	Id          string `json:"id"`
	PublishedAt string `json:"publishedAt"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
}
