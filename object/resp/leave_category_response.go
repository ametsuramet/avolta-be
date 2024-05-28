package resp

type LeaveCategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Absent      bool   `json:"absent"`
	Sick        bool   `json:"sick"`
}
