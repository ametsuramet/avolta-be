package model

import (
	"avolta/object/resp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reimbursement struct {
	Base
	AccountPayableID string              `gorm:"size:30" json:"account_payable_id"`
	AccountPayable   Account             `gorm:"foreignKey:AccountPayableID" json:"-"`
	AccountExpenseID string              `gorm:"size:30" json:"account_expense_id"`
	AccountExpense   Account             `gorm:"foreignKey:AccountExpenseID" json:"-"`
	Date             time.Time           `json:"date" binding:"required"`
	Notes            string              `json:"notes"`
	Total            float64             `json:"total"`
	Balance          float64             `json:"balance"`
	Status           string              `json:"status" gorm:"type:enum('DRAFT', 'REQUEST', 'PROCESSING', 'APPROVED', 'REJECTED', 'PAID', 'FINISHED', 'CANCELED');default:'DRAFT'"`
	Items            []ReimbursementItem `json:"item,omitempty"`
	EmployeeID       string              `json:"employee_id"`
	Employee         Employee            `gorm:"foreignKey:EmployeeID" json:"-"`
	Transactions     []Transaction       `json:"transactions,omitempty"`
	Attachment       string              `json:"attachment"`
}

func (u *Reimbursement) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Reimbursement) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ReimbursementReponse{})
}
