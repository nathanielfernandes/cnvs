package preview

type PreviewResponse struct {
	AudioURL    string `json:"audio_url"`
	CoverArtURL string `json:"cover_art_url"`
	TrackName   string `json:"track_name"`
	ArtistName  string `json:"artist_name"`
	AlbumName   string `json:"album_name"`
}
