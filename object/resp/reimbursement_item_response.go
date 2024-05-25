package resp

type ReimbursementItemResponse struct {
	ID          string   `json:"id" `
	Amount      float64  `json:"amount" `
	Notes       string   `json:"notes" `
	Attachments []string `json:"attachments" gorm:"-"`
}
