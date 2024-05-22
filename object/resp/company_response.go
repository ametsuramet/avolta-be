package resp

type CompanyReponse struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Logo                  string `json:"logo"`
	Cover                 string `json:"cover"`
	LegalEntity           string `json:"legal_entity"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	Fax                   string `json:"fax"`
	Address               string `json:"address"`
	ContactPerson         string `json:"contact_person"`
	ContactPersonPosition string `json:"contact_person_position"`
	TaxPayerNumber        string `json:"tax_payer_number"`
	LogoURL               string `json:"logo_url"`
	CoverURL              string `json:"cover_url"`
}
