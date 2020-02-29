package types

type MovieData struct {
	Year      int64        `json:"year"`
	Title     string       `json:"title"`
	ID        int64        `json:"tmdbId"`
	Quality   int64        `json:"qualityProfileId"`
	TitleSlug string       `json:"titleSlug"`
	Images    []Image      `json:"images"`
	Path      string       `json:"rootFolderPath"`
	Monitored bool         `json:"monitored"`
	Options   MovieOptions `json:"addOptions"`
}

type Image struct {
	CoverType string `json:"coverType"`
	Url       string `json:"url"`
}

type MovieOptions struct {
	Search bool `json:"searchForMovie"`
}
