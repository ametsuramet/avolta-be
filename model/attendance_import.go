package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceBulkImport struct {
	Base
	FileName       string             `json:"file_name"`
	ImportedBy     string             `json:"-"`
	User           User               `json:"user" gorm:"foreignKey:ImportedBy"`
	DateImportedAt time.Time          `json:"date_imported_at"`
	Data           []AttendanceImport `json:"data" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status         string             `json:"status" gorm:"type:enum('DRAFT','APPROVED','REJECTED') DEFAULT 'DRAFT'"`
	Notes          string
	CompanyID      string  `json:"company_id" gorm:"not null"`
	Company        Company `gorm:"foreignKey:CompanyID"`
}
type AttendanceImport struct {
	Base
	SequenceNumber         int                    `json:"sequence_number"`
	FingerprintID          string                 `json:"fingerprint_id"`
	EmployeeCode           string                 `json:"employee_code"`
	EmployeeName           string                 `json:"employee_name"`
	SystemEmployeeName     string                 `json:"system_employee_name"`
	SystemEmployeeID       string                 `json:"system_employee_id"`
	SystemEmployee         Employee               `json:"-" gorm:"foreignKey:SystemEmployeeID"`
	Items                  []AttendanceImportItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AttendanceBulkImportID string                 `json:"-"`
	AttendanceBulkImport   AttendanceBulkImport   `json:"-" gorm:"foreignKey:AttendanceBulkImportID"`
}

type AttendanceImportItem struct {
	Base
	SequenceNumber     int              `json:"sequence_number"`
	Day                string           `json:"day"`
	Date               string           `json:"date"`
	WorkingHour        string           `json:"working_hour"`
	Activity           string           `json:"activity"`
	DutyOn             string           `json:"duty_on"`
	DutyOff            string           `json:"duty_off"`
	LateIn             string           `json:"late_in"`
	EarlyDeparture     string           `json:"early_departure"`
	EffectiveHour      string           `json:"effective_hour"`
	Overtime           string           `json:"overtime"`
	Notes              string           `json:"notes"`
	AttendanceImportID string           `json:"-"`
	AttendanceImport   AttendanceImport `json:"-" gorm:"foreignKey:AttendanceImportID"`
}

func (u *AttendanceBulkImport) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}
func (u *AttendanceImport) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}
func (u *AttendanceImportItem) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}
