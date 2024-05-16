package resp

type UserResponse struct {
	FullName    string   `json:"full_name"`
	Picture     string   `json:"picture"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	RoleName    string   `json:"role_name"`
}
