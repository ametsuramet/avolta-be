package model

import (
	"avolta/config"
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	Base
	Parent           *Organization
	ParentId         *string        `gorm:"REFERENCES organizations" json:"parent_id"`
	Name             string         `gorm:"size:30" json:"name"`
	Code             string         `gorm:"size:30" json:"code"`
	Description      string         `gorm:"size:100" json:"description"`
	Employee         []Employee     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SubOrganizations []Organization `json:"sub_organizations" gorm:"foreignKey:parent_id"`
}

func (u *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Organization) MarshalJSON() ([]byte, error) {
	parent := ""
	parentId := ""
	if m.Parent != nil {
		parent = m.Parent.Name
		parentId = *m.ParentId
	}
	employees := []resp.SimpleEmployeeResponse{}
	// subOrganizations := []resp.OrganizationResponse{}

	for _, v := range m.Employee {
		var pictureUrl string
		if v.Picture.Valid {
			pictureUrl = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, v.Picture.String)
		}
		employees = append(employees, resp.SimpleEmployeeResponse{
			ID:         v.ID,
			FullName:   v.FullName,
			JobTitle:   v.JobTitle.Name,
			PictureUrl: pictureUrl,
		})
	}

	subOrganizations := m.GetSuborganizationResponse()

	return json.Marshal(resp.OrganizationResponse{
		ID:               m.ID,
		Parent:           parent,
		ParentId:         parentId,
		Name:             m.Name,
		Code:             m.Code,
		Description:      m.Description,
		Employees:        employees,
		SubOrganizations: subOrganizations,
	})
}

func (m *Organization) GetSuborganizationResponse() []resp.OrganizationResponse {
	m.SubOrganizations = m.GetSuborganization()
	subOrganizations := []resp.OrganizationResponse{}
	for _, v := range m.SubOrganizations {

		parentName := m.Name
		parentId := m.ID
		// if v.Parent != nil {
		// 	parentName = v.Parent.Name
		// 	parentId = v.Parent.ID
		// }

		employees := []resp.SimpleEmployeeResponse{}

		for _, v := range m.Employee {
			var pictureUrl string
			if v.Picture.Valid {
				pictureUrl = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, v.Picture.String)
			}
			employees = append(employees, resp.SimpleEmployeeResponse{
				ID:         v.ID,
				FullName:   v.FullName,
				JobTitle:   v.JobTitle.Name,
				PictureUrl: pictureUrl,
			})
		}

		subOrganizations = append(subOrganizations, resp.OrganizationResponse{
			ID:               v.ID,
			Parent:           parentName,
			ParentId:         parentId,
			Name:             v.Name,
			Code:             v.Code,
			Description:      v.Description,
			Employees:        employees,
			SubOrganizations: v.GetSuborganizationResponse(),
		})
	}

	return subOrganizations
}

func (m *Organization) GetSuborganization() []Organization {
	suborgs := []Organization{}
	database.DB.Find(&m.SubOrganizations, "parent_id = ?", m.ID)
	for _, v := range m.SubOrganizations {
		suborgs = append(suborgs, Organization{
			Base:             Base{ID: v.ID},
			Parent:           m,
			ParentId:         &m.ID,
			Name:             v.Name,
			Code:             v.Code,
			Description:      v.Description,
			Employee:         v.Employee,
			SubOrganizations: v.GetSuborganization(),
		})
	}

	return suborgs
}
