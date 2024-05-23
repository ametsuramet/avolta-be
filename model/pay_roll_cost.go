package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRollCost struct {
	Base
	Description   string
	PayRollID     string      `json:"pay_roll_id"`
	PayRoll       PayRoll     `gorm:"foreignKey:PayRollID" json:"-"`
	PayRollItemID string      `json:"pay_roll_item_id"`
	PayRollItem   PayRollItem `gorm:"foreignKey:PayRollItemID" json:"-"`
	Amount        float64     `json:"amount"`
	Tariff        float64     `json:"tariff"`
	BpjsTkJht     bool        `json:"bpjs_tk_jht"`
	BpjsTkJp      bool        `json:"bpjs_tk_jp"`
	DebtDeposit   bool        `json:"debt_deposit"`
}

func (u *PayRollCost) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m PayRollCost) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.PayRollCostReponse{})
}
