package resp

type BlogPostReponse struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	AuthorID   string          `json:"author_id"`
	CategoryID string          `json:"category_id"`
	Published  bool            `json:"published"`
	Category   string          `json:"category"`
	Author     string          `json:"author"`
	Tags       []*TagReponse   `json:"tags"`
	Images     []*ImageReponse `json:"images"`
	IsPage     bool            `json:"is_page"`
}
