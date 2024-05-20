package resp

type OrganizationReponse struct {
	ID               string                  `json:"id"`
	Parent           string                  `json:"parent"`
	ParentId         string                  `json:"parent_id"`
	Name             string                  `json:"name"`
	Code             string                  `json:"code"`
	Description      string                  `json:"description"`
	Employees        []SimpleEmployeeReponse `json:"employees"`
	SubOrganizations []OrganizationReponse   `json:"sub_organizations"`
}
