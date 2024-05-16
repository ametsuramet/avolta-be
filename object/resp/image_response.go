package resp

type ImageReponse struct {
	Path    string `json:"path"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}
