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

type Leave struct {
	Base
	Name            string        `json:"name"`
	RequestType     string        `json:"request_type" gorm:"type:enum('FULL_DAY', 'HALF_DAY', 'HOURLY');default:'FULL_DAY'"`
	LeaveCategoryID *string       `json:"leave_category_id"`
	LeaveCategory   LeaveCategory `json:"leave_category" gorm:"foreignKey:LeaveCategoryID"`
	StartDate       *time.Time    `json:"start_date" gorm:"type:DATE" sql:"TYPE:DATE"`
	EndDate         *time.Time    `json:"end_date" gorm:"type:DATE" sql:"TYPE:DATE"`
	StartTime       *TimeOnly     `json:"start_time" gorm:"type:TIME"`
	EndTime         *TimeOnly     `json:"end_time" gorm:"type:TIME"`
	EmployeeID      string        `json:"employee_id"`
	Employee        Employee      `json:"employee" gorm:"foreignKey:EmployeeID"`
	Description     string        `json:"description" gorm:"type:TEXT"`
	Status          string        `json:"status" gorm:"type:enum('DRAFT','REVIEWED','APPROVED','REJECTED') DEFAULT 'DRAFT'"`
	Remarks         string        `json:"remarks" gorm:"type:TEXT"`
	Attachment      *string       `json:"attachment" gorm:"type:TEXT"`
	ApproverID      *string       `json:"approver_id"`
	Approver        User          `json:"approver" gorm:"foreignKey:ApproverID"`
	CompanyID       string        `json:"company_id" gorm:"not null"`
	Company         Company       `gorm:"foreignKey:CompanyID"`
}

func (u *Leave) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Leave) MarshalJSON() ([]byte, error) {

	startDate := ""
	endDate := ""
	startTime := ""
	endTime := ""
	if m.StartDate != nil {
		startDate = m.StartDate.Format("2006-01-02")
	}

	if m.EndDate != nil {
		endDate = m.EndDate.Format("2006-01-02")
	}
	if m.StartTime != nil {
		startTime = m.StartTime.Format("15:04:05")
	}

	if m.EndTime != nil {
		endTime = m.EndTime.Format("15:04:05")
	}

	var attachmentURL string
	if m.Attachment != nil {
		attachmentURL = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, *m.Attachment)
	}

	var employeePicture string
	if m.Employee.Picture.Valid {
		employeePicture = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Employee.Picture.String)
	}
	diffDays := int64(0)
	if m.RequestType == "FULL_DAY" {
		diff := m.EndDate.Sub(*m.StartDate)
		diffDays = int64(diff.Hours()/24) + 1
	}

	return json.Marshal(resp.LeaveResponse{
		ID:              m.ID,
		Name:            m.Name,
		RequestType:     m.RequestType,
		LeaveCategoryID: m.LeaveCategory.ID,
		LeaveCategory:   m.LeaveCategory.Name,
		StartDate:       startDate,
		EndDate:         endDate,
		StartTime:       startTime,
		EndTime:         endTime,
		EmployeeID:      m.EmployeeID,
		EmployeeName:    m.Employee.FullName,
		EmployeePicture: employeePicture,
		Description:     m.Description,
		Status:          m.Status,
		Remarks:         m.Remarks,
		AttachmentURL:   attachmentURL,
		Absent:          m.LeaveCategory.Absent,
		Sick:            m.LeaveCategory.Sick,
		Diff:            diffDays,
	})
}
