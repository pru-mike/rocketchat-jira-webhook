package confluence

type PageSearchResult struct {
	Size    int `json:"size"`
	Start   int `json:"start"`
	Limit   int `json:"limit"`
	Results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"results"`
}
