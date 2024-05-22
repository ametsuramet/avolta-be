package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	Base
	PayRollAutoNumber                bool    `json:"pay_roll_auto_number" gorm:"default:true"`
	PayRollAutoFormat                string  `json:"pay_roll_auto_format" gorm:"default:'{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}'"`
	PayRollStaticCharacter           string  `json:"pay_roll_static_character" gorm:"default:'PAYROLL'"`
	PayRollAutoNumberCharacterLength int     `json:"pay_roll_auto_number_character_length" gorm:"default:5"`
	PayRollPayableAccountID          *string `json:"pay_roll_payable_account_id"`
	PayRollPayableAccount            Account `gorm:"foreignKey:PayRollPayableAccountID" json:"pay_roll_payable_account"`
	PayRollExpenseAccountID          *string `json:"pay_roll_expense_account_id"`
	PayRollExpenseAccount            Account `gorm:"foreignKey:PayRollExpenseAccountID" json:"pay_roll_expense_account"`
	PayRollAssetAccountID            *string `json:"pay_roll_asset_account_id"`
	PayRollAssetAccount              Account `gorm:"foreignKey:PayRollAssetAccountID" json:"pay_roll_asset_account"`
	PayRollTaxAccountID              *string `json:"pay_roll_tax_account_id"`
	PayRollTaxAccount                Account `gorm:"foreignKey:PayRollTaxAccountID" json:"pay_roll_tax_account"`
	IsEffectiveRateAverage           bool    `json:"is_effective_rate_average"`
	IsGrossUp                        bool    `json:"is_gross_up"`
}

func (u *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Setting) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.SettingReponse{
		ID:                               m.ID,
		PayRollAutoNumber:                m.PayRollAutoNumber,
		PayRollAutoFormat:                m.PayRollAutoFormat,
		PayRollStaticCharacter:           m.PayRollStaticCharacter,
		PayRollAutoNumberCharacterLength: m.PayRollAutoNumberCharacterLength,
		PayRollPayableAccountID:          m.PayRollPayableAccountID,
		PayRollExpenseAccountID:          m.PayRollExpenseAccountID,
		PayRollAssetAccountID:            m.PayRollAssetAccountID,
		PayRollTaxAccountID:              m.PayRollTaxAccountID,
		IsEffectiveRateAverage:           m.IsEffectiveRateAverage,
		IsGrossUp:                        m.IsGrossUp,
	})
}
