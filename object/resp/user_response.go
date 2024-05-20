package resp

type UserResponse struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name"`
	Picture     string   `json:"picture"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	RoleName    string   `json:"role_name"`
	RoleID      string   `json:"role_id"`
	IsAdmin     bool     `json:"is_admin"`
}
