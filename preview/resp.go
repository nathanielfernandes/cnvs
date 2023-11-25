package preview

type PreviewResponse struct {
	AudioURL        string   `json:"audio_url"`
	CoverArt        CoverArt `json:"cover_art"`
	TrackName       string   `json:"track_name"`
	Artists         []Artist `json:"artists"`
	AlbumName       string   `json:"album_name,omitempty"`
	BackgroundColor string   `json:"background_color"`
	ReleaseDate     string   `json:"release_date,omitempty"`
}

type Artist struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CoverArt struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}
