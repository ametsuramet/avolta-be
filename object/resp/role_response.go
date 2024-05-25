package resp

type RoleResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	Permissions  []string `json:"permissions"`
}
