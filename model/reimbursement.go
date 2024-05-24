package model

import (
	"avolta/config"
	"avolta/object/resp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reimbursement struct {
	Base
	AccountPayableID *string             `gorm:"size:30" json:"account_payable_id"`
	AccountPayable   Account             `gorm:"foreignKey:AccountPayableID" json:"-"`
	AccountExpenseID *string             `gorm:"size:30" json:"account_expense_id"`
	AccountExpense   Account             `gorm:"foreignKey:AccountExpenseID" json:"-"`
	Date             time.Time           `json:"date" binding:"required"`
	Name             string              `json:"name"`
	Notes            string              `json:"notes"`
	Remarks          string              `json:"remarks"`
	Total            float64             `json:"total"`
	Balance          float64             `json:"balance"`
	Status           string              `json:"status" gorm:"type:enum('DRAFT', 'REQUEST', 'PROCESSING', 'APPROVED', 'REJECTED', 'PAID', 'FINISHED', 'CANCELED');default:'DRAFT'"`
	Items            []ReimbursementItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EmployeeID       string              `json:"employee_id"`
	Employee         Employee            `gorm:"foreignKey:EmployeeID" json:"-"`
	Transactions     []Transaction       `json:"transactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachment       string              `json:"attachment"`
}

func (u *Reimbursement) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Reimbursement) MarshalJSON() ([]byte, error) {
	files := []string{}
	attachments := []string{}

	if m.Attachment != "" {
		json.Unmarshal([]byte(m.Attachment), &files)
		for _, file := range files {
			attachments = append(attachments, fmt.Sprintf("%s/%s", config.App.Server.BaseURL, file))
		}
	}

	return json.Marshal(resp.ReimbursementReponse{
		ID:      m.ID,
		Date:    m.Date.Format(time.RFC3339),
		Notes:   m.Notes,
		Remarks: m.Remarks,
		Name:    m.Name,
		Total:   m.Total,
		Balance: m.Balance,
		Status:  m.Status,
		Items: func(items []ReimbursementItem) []resp.ReimbursementItemReponse {
			var newItems = []resp.ReimbursementItemReponse{}
			for _, v := range items {
				attachments := []string{}
				files := []string{}
				json.Unmarshal([]byte(v.Files), &files)
				for _, file := range files {
					attachments = append(attachments, fmt.Sprintf("%s/%s", config.App.Server.BaseURL, file))
				}
				newItems = append(newItems, resp.ReimbursementItemReponse{ID: v.ID, Amount: v.Amount, Notes: v.Notes, Attachments: attachments})
			}
			return newItems
		}(m.Items),
		EmployeeID:   m.EmployeeID,
		EmployeeName: m.Employee.FullName,
		Transactions: func(items []Transaction) []resp.TransactionReponse {
			var newItems = []resp.TransactionReponse{}
			for _, v := range items {
				accountSourceName := ""
				if v.AccountSourceID != nil {
					accountSourceName = v.AccountSource.Name
				}
				accountDestinationName := ""
				if v.AccountDestinationID != nil {
					accountDestinationName = v.AccountDestination.Name
				}
				newItems = append(newItems, resp.TransactionReponse{ID: v.ID, Date: v.Date.Format(time.RFC3339), Amount: v.Amount, Credit: v.Credit, Debit: v.Debit, Notes: v.Notes, Description: v.Description, AccountSourceName: accountSourceName, AccountDestinationName: accountDestinationName})
			}
			return newItems
		}(m.Transactions),
		Attachments: attachments,
	})
}
