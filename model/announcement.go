package model

import (
	"avolta/object/resp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Announcement struct {
	Base
	Name           string
	Description    string
	StartDate      time.Time
	EndDate        *time.Time
	OrganizationID *string      `json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	CompanyID      string       `json:"company_id" gorm:"not null"`
	Company        Company      `gorm:"foreignKey:CompanyID"`
}

func (u *Announcement) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Announcement) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.AnnouncementResponse{})
}
