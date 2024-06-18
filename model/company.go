package model

import (
	"avolta/config"
	"avolta/object/resp"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	Base
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
	UserID                string `json:"user_id"`
}

func (u *Company) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Company) MarshalJSON() ([]byte, error) {
	var coverURL string
	if m.Cover != "" {
		coverURL = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Cover)
	}
	var logoURL string
	if m.Logo != "" {
		logoURL = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Logo)
	}
	return json.Marshal(resp.CompanyResponse{
		ID:                    m.ID,
		Name:                  m.Name,
		Logo:                  m.Logo,
		Cover:                 m.Cover,
		LegalEntity:           m.LegalEntity,
		Email:                 m.Email,
		Phone:                 m.Phone,
		Fax:                   m.Fax,
		Address:               m.Address,
		ContactPerson:         m.ContactPerson,
		ContactPersonPosition: m.ContactPersonPosition,
		TaxPayerNumber:        m.TaxPayerNumber,
		LogoURL:               logoURL,
		CoverURL:              coverURL,
	})
}
