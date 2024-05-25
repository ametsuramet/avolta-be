package resp

type ImageResponse struct {
	Path    string `json:"path"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}
