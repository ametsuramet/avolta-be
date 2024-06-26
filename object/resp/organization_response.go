package resp

type OrganizationResponse struct {
	ID               string                   `json:"id"`
	Parent           string                   `json:"parent"`
	ParentId         string                   `json:"parent_id"`
	Name             string                   `json:"name"`
	Code             string                   `json:"code"`
	Description      string                   `json:"description"`
	Employees        []SimpleEmployeeResponse `json:"employees"`
	SubOrganizations []OrganizationResponse   `json:"sub_organizations"`
}
