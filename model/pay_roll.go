package model

import (
	"avolta/config"
	"avolta/database"
	"avolta/object/resp"
	"avolta/service"
	"avolta/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/apung/go-terbilang"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leekchan/accounting"
	"gorm.io/gorm"
)

type PayRoll struct {
	Base
	PayRollNumber                   string          `json:"pay_roll_number"`
	Title                           string          `json:"title"`
	Notes                           string          `json:"notes"`
	StartDate                       time.Time       `json:"start_date" binding:"required"`
	EndDate                         time.Time       `json:"end_date" binding:"required"`
	Files                           string          `json:"files" gorm:"default:'[]'"`
	TotalIncome                     float64         `json:"total_income"`
	TotalReimbursement              float64         `json:"total_reimbursement"`
	TotalDeduction                  float64         `json:"total_deduction"`
	TotalTax                        float64         `json:"total_tax"`
	TaxCost                         float64         `json:"tax_cost"`
	NetIncome                       float64         `json:"net_income"`
	NetIncomeBeforeTaxCost          float64         `json:"net_income_before_tax_cost"`
	TakeHomePay                     float64         `json:"take_home_pay"`
	TotalPayable                    float64         `json:"total_payable"`
	TaxAllowance                    float64         `json:"tax_allowance"`
	TaxTariff                       float64         `json:"tax_tariff"`
	IsGrossUp                       bool            `json:"is_gross_up"`
	IsEffectiveRateAverage          bool            `json:"is_effective_rate_average"`
	Status                          string          `json:"status" gorm:"type:enum('DRAFT', 'RUNNING', 'FINISHED');default:'DRAFT'"`
	Attachments                     []string        `json:"attachments" gorm:"-"`
	Transactions                    []Transaction   `json:"transactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PayableTransactions             []Transaction   `json:"payable_transactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items                           []PayRollItem   `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Costs                           []PayRollCost   `json:"costs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TakeHomePayCounted              string          `json:"take_home_pay_counted" gorm:"-"`
	TakeHomePayReimbursementCounted string          `json:"take_home_pay_reimbursement_counted" gorm:"-"`
	TaxPaymentID                    sql.NullString  `json:"tax_payment_id"`
	EmployeeID                      string          `binding:"required" json:"employee_id"`
	Employee                        Employee        `gorm:"foreignKey:EmployeeID" `
	TaxSummary                      CountTaxSummary `gorm:"-" json:"tax_summary"`
	BpjsSetting                     *service.Bpjs   `gorm:"-" `
	PayRollReportItemID             *string         `json:"pay_roll_report_item_id"`
	CompanyID                       string          `json:"company_id" gorm:"not null"`
	Company                         Company         `gorm:"foreignKey:CompanyID"`
}

type PayRollPayment struct {
	TransactionRefID     string    `json:"transaction_ref_id"`
	Date                 time.Time `json:"date"`
	Amount               float64   `json:"amount"`
	AccountDestinationID string    `json:"account_destination_id"`
	Description          string    `json:"description"`
}

func (u *PayRoll) BeforeCreate(tx *gorm.DB) (err error) {

	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m PayRoll) MapWithOutDetail() resp.PayRollResponse {
	return resp.PayRollResponse{
		ID:                              m.ID,
		PayRollNumber:                   m.PayRollNumber,
		Notes:                           m.Notes,
		StartDate:                       m.StartDate.Format("2006-01-02"),
		EndDate:                         m.EndDate.Format("2006-01-02"),
		EmployeeID:                      m.EmployeeID,
		TaxAllowance:                    m.TaxAllowance,
		EmployeeName:                    m.Employee.FullName,
		TaxSummary:                      resp.CountTaxSummary(m.TaxSummary),
		IsGrossUp:                       m.IsGrossUp,
		IsEffectiveRateAverage:          m.IsEffectiveRateAverage,
		TotalIncome:                     m.TotalIncome,
		TotalReimbursement:              m.TotalReimbursement,
		TotalDeduction:                  m.TotalDeduction,
		TotalTax:                        m.TotalTax,
		TaxCost:                         m.TaxCost,
		NetIncome:                       m.NetIncome,
		NetIncomeBeforeTaxCost:          m.NetIncomeBeforeTaxCost,
		TakeHomePay:                     m.TakeHomePay,
		TotalPayable:                    m.TotalPayable,
		TakeHomePayCounted:              m.TakeHomePayCounted,
		TakeHomePayReimbursementCounted: m.TakeHomePayReimbursementCounted,
		Status:                          m.Status,
	}
}
func (m PayRoll) MarshalJSON() ([]byte, error) {

	m.GetItems()

	items := []resp.PayRollItemResponse{}
	transactions := []resp.TransactionResponse{}

	for _, v := range m.Items {
		items = append(items, resp.PayRollItemResponse{
			ID:             v.ID,
			ItemType:       v.ItemType,
			Title:          v.Title,
			Notes:          v.Notes,
			IsDefault:      v.IsDefault,
			IsDeductible:   v.IsDeductible,
			IsTax:          v.IsTax,
			TaxAutoCount:   v.TaxAutoCount,
			IsTaxCost:      v.IsTaxCost,
			IsTaxAllowance: v.IsTaxAllowance,
			Amount:         v.Amount,
			Tariff:         v.Tariff,
		})
	}

	m.GetPayRollCost()
	costs := []resp.PayRollCostResponse{}
	for _, v := range m.Costs {
		costs = append(costs, resp.PayRollCostResponse{
			ID:          v.ID,
			Description: v.Description,
			Amount:      v.Amount,
			BpjsTkJht:   v.BpjsTkJht,
			BpjsTkJp:    v.BpjsTkJp,
			Tariff:      v.Tariff,
			DebtDeposit: v.DebtDeposit,
		})
	}

	for _, v := range m.Transactions {
		accountSourceId := ""
		accountSourceName := ""
		if v.AccountSourceID != nil {
			accountSourceId = *v.AccountSourceID
			accountSourceName = v.AccountSource.Name
		}
		accountDestinationId := ""
		accountDestinationName := ""
		if v.AccountDestinationID != nil {
			accountDestinationId = *v.AccountDestinationID
			accountDestinationName = v.AccountDestination.Name
		}
		transactions = append(transactions, resp.TransactionResponse{
			ID:                     v.ID,
			Description:            v.Description,
			Notes:                  v.Notes,
			Credit:                 v.Credit,
			Debit:                  v.Debit,
			Amount:                 v.Amount,
			Date:                   v.Date.Format("2006-01-02 15:04:05"),
			IsIncome:               v.IsIncome,
			IsExpense:              v.IsExpense,
			IsJournal:              v.IsJournal,
			IsAccountReceivable:    v.IsAccountReceivable,
			IsAccountPayable:       v.IsAccountPayable,
			AccountSourceID:        accountSourceId,
			AccountSourceName:      accountSourceName,
			AccountDestinationID:   accountDestinationId,
			AccountDestinationName: accountDestinationName,
			EmployeeID:             v.EmployeeID,
			EmployeeName:           v.Employee.FullName,
			IsPayRollPayment:       v.IsPayRollPayment,
			IsReimbursementPayment: v.IsReimbursementPayment,
		})
	}

	payableTransactions := []resp.TransactionResponse{}
	for _, v := range m.PayableTransactions {
		accountSourceId := ""
		accountSourceName := ""
		if v.AccountSourceID != nil {
			accountSourceId = *v.AccountSourceID
			accountSourceName = v.AccountSource.Name
		}
		accountDestinationId := ""
		accountDestinationName := ""
		if v.AccountDestinationID != nil {
			accountDestinationId = *v.AccountDestinationID
			accountDestinationName = v.AccountDestination.Name
		}
		payableTransactions = append(payableTransactions, resp.TransactionResponse{
			ID:                     v.ID,
			Description:            v.Description,
			Notes:                  v.Notes,
			Credit:                 v.Credit,
			Debit:                  v.Debit,
			Amount:                 v.Amount,
			Date:                   v.Date.Format("2006-01-02 15:04:05"),
			IsIncome:               v.IsIncome,
			IsExpense:              v.IsExpense,
			IsJournal:              v.IsJournal,
			IsAccountReceivable:    v.IsAccountReceivable,
			IsAccountPayable:       v.IsAccountPayable,
			AccountSourceID:        accountSourceId,
			AccountSourceName:      accountSourceName,
			AccountDestinationID:   accountDestinationId,
			AccountDestinationName: accountDestinationName,
			EmployeeID:             v.EmployeeID,
			EmployeeName:           v.Employee.FullName,
			IsPayRollPayment:       v.IsPayRollPayment,
			IsReimbursementPayment: v.IsReimbursementPayment,
		})
	}

	m.TakeHomePayCounted = strings.ToTitle(strings.ToLower(terbilang.ToTerbilang(int(math.Round(m.TakeHomePay)))))
	m.TakeHomePayReimbursementCounted = strings.ToTitle(strings.ToLower(terbilang.ToTerbilang(int(math.Round(m.TakeHomePay + m.TotalReimbursement)))))
	response := m.MapWithOutDetail()
	response.Items = items
	response.Transactions = transactions
	response.PayableTransactions = payableTransactions
	response.Costs = costs
	return json.Marshal(response)
}

func (m *PayRoll) GetEmployee() {
	employee := Employee{}
	database.DB.Find(&employee, "id = ?", m.EmployeeID)
	m.Employee = employee
}

func (m *PayRoll) Payment(c *gin.Context) error {
	var data PayRollPayment

	if err := c.ShouldBindJSON(&data); err != nil {
		return err
	}
	var setting Setting
	if err := database.DB.First(&setting).Error; err != nil {
		return err
	}
	refTrans := Transaction{}
	database.DB.Find(&refTrans, "id = ?", data.TransactionRefID)
	if err := database.DB.Create(&Transaction{
		Description:          data.Description + " (" + m.PayRollNumber + ")",
		Debit:                data.Amount,
		AccountDestinationID: setting.PayRollPayableAccountID,
		Date:                 data.Date,
		PayRollID:            &m.ID,
		EmployeeID:           m.EmployeeID,
		IsPayRollPayment:     true,
		TransactionRefID:     &data.TransactionRefID,
		IsTakeHomePay:        refTrans.IsTakeHomePay,
	}).Error; err != nil {
		return err
	}
	return nil
}
func (m *PayRoll) RunPayRoll(c *gin.Context) error {
	m.GetItems()

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var setting Setting
		if err := database.DB.First(&setting).Error; err != nil {
			return err
		}

		now := time.Now()

		// GET REIMBURSEMENT
		reimbursementItem := PayRollItem{}
		tx.Find(&reimbursementItem, "pay_roll_id = ? AND reimbursement_id is not null", m.ID)
		_, totalReimbursement, _, _ := m.GetDeductible()

		m.TotalReimbursement = totalReimbursement

		// CREATE TRANSACTION
		// TAKE HOME PAY
		if err := tx.Create(&Transaction{

			Description:          "Take Home Pay (" + m.PayRollNumber + ")",
			Debit:                m.TakeHomePay,
			AccountDestinationID: setting.PayRollExpenseAccountID,
			IsExpense:            true,
			Date:                 now,
			PayRollID:            &m.ID,
			EmployeeID:           m.EmployeeID,
			IsTakeHomePay:        true,
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&Transaction{

			Description:          "Hutang Take Home Pay (" + m.PayRollNumber + ")",
			Credit:               m.TakeHomePay,
			AccountDestinationID: setting.PayRollPayableAccountID,
			IsAccountPayable:     true,
			Date:                 now,
			PayRollID:            &m.ID,
			EmployeeID:           m.EmployeeID,
			IsTakeHomePay:        true,
		}).Error; err != nil {
			return err
		}

		// REIMBURSEMENT
		fmt.Println("TOTAL REIMBURSEMENT", totalReimbursement)
		if totalReimbursement > 0 {
			reimbursementItems := m.GetReimbursementItems()

			for _, v := range reimbursementItems {
				if v.ReimbursementID != nil {
					var reimbursement Reimbursement
					if err := tx.Find(&reimbursement, "id = ?", &v.ReimbursementID).Error; err != nil {
						return err
					}
					reimbursement.Status = "FINISHED"
					reimbursement.Notes = "[MOVED TO PAYROLL] -> " + m.PayRollNumber
					if err := tx.Model(&reimbursement).Updates(&reimbursement).Error; err != nil {
						return err
					}

					// if err := tx.Create(&Transaction{

					// 	Description:          "MOVE TO PAYROLL (" + reimbursement.Name + " -> " + m.PayRollNumber + ")",
					// 	Credit:               reimbursement.Total,
					// 	AccountDestinationID: reimbursement.AccountExpenseID,
					// 	IsExpense:            true,
					// 	Date:                 now,
					// 	ReimbursementID:      &reimbursement.ID,
					// 	EmployeeID:           reimbursement.EmployeeID,
					// 	PayRollID:            &m.ID,
					// }).Error; err != nil {
					// 	return err
					// }
					// if err := tx.Create(&Transaction{

					// 	Description:          "MOVE TO PAYROLL (" + reimbursement.Name + " -> " + m.PayRollNumber + ")",
					// 	Debit:                reimbursement.Total,
					// 	AccountDestinationID: reimbursement.AccountPayableID,
					// 	IsAccountPayable:     true,
					// 	Date:                 now,
					// 	ReimbursementID:      &reimbursement.ID,
					// 	EmployeeID:           reimbursement.EmployeeID,
					// 	PayRollID:            &m.ID,
					// }).Error; err != nil {
					// 	return err
					// }
				}
			}

			if err := tx.Create(&Transaction{

				Description:          "Reimbursement (" + m.PayRollNumber + ")",
				Debit:                m.TotalReimbursement,
				AccountDestinationID: setting.PayRollExpenseAccountID,
				IsExpense:            true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&Transaction{

				Description:          "Hutang Reimbursement (" + m.PayRollNumber + ")",
				Credit:               m.TotalReimbursement,
				AccountDestinationID: setting.PayRollPayableAccountID,
				IsAccountPayable:     true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}
		}

		if m.TotalTax > 0 {
			//  TAX TRANSACTION
			if err := tx.Create(&Transaction{

				Description:          "Pajak (" + m.PayRollNumber + ")",
				Debit:                m.TotalTax,
				AccountDestinationID: setting.PayRollExpenseAccountID,
				IsExpense:            true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&Transaction{

				Description:          "Hutang Pajak (" + m.PayRollNumber + ")",
				Credit:               m.TotalTax,
				AccountDestinationID: setting.PayRollTaxAccountID,
				IsAccountPayable:     true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}
		}
		m.GetPayRollCost()

		for _, v := range m.Costs {
			if err := tx.Create(&Transaction{
				Description:          v.Description + " (" + m.PayRollNumber + ")",
				Debit:                v.Amount,
				AccountDestinationID: setting.PayRollExpenseAccountID,
				IsExpense:            true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&Transaction{
				Description:          "Hutang " + v.Description + " (" + m.PayRollNumber + ")",
				Credit:               v.Amount,
				AccountDestinationID: setting.PayRollCostAccountID,
				IsAccountPayable:     true,
				Date:                 now,
				PayRollID:            &m.ID,
				EmployeeID:           m.EmployeeID,
			}).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&m).Update("status", "RUNNING").Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
func (m *PayRoll) CreateDefaultItems(c *gin.Context) error {
	m.GetEmployee()

	//  CREATE BASIC SALARY
	if err := database.DB.Create(&PayRollItem{
		ItemType:    config.SALARY,
		IsDefault:   true,
		Amount:      m.Employee.BasicSalary,
		PayRollID:   m.ID,
		Title:       "Gaji Pokok",
		BpjsCounted: true,
	}).Error; err != nil {
		return err
	}

	//  CREATE POSITIONAL ALLOWANCE
	if err := database.DB.Create(&PayRollItem{
		ItemType:    config.ALLOWANCE,
		IsDefault:   false,
		Amount:      m.Employee.PositionalAllowance,
		PayRollID:   m.ID,
		Title:       "Tunjangan Jabatan",
		BpjsCounted: true,
	}).Error; err != nil {
		return err
	}

	//  CREATE TRANSPORT ALLOWANCE
	if err := database.DB.Create(&PayRollItem{
		ItemType:  config.ALLOWANCE,
		IsDefault: false,
		Amount:    m.Employee.TransportAllowance,
		PayRollID: m.ID,
		Title:     "Uang Transport",
	}).Error; err != nil {
		return err
	}
	//  CREATE MEAL ALLOWANCE
	if err := database.DB.Create(&PayRollItem{
		ItemType:  config.ALLOWANCE,
		IsDefault: false,
		Amount:    m.Employee.MealAllowance,
		PayRollID: m.ID,
		Title:     "Uang Makan",
	}).Error; err != nil {
		return err
	}

	// CREATE OVERTIME
	if err := database.DB.Create(&PayRollItem{
		ItemType:  config.OVERTIME,
		IsDefault: false,
		Amount:    0,
		PayRollID: m.ID,
		Title:     "Lembur",
	}).Error; err != nil {
		return err
	}

	// CREATE REIMBURSEMENT
	if err := database.DB.Create(&PayRollItem{
		ItemType:  config.REIMBURSEMENT,
		IsDefault: false,
		Amount:    0,
		PayRollID: m.ID,
		Title:     "Reimbursement",
	}).Error; err != nil {
		return err
	}

	// DEDUCTION

	//  POSITIONAL COST
	if err := database.DB.Create(&PayRollItem{
		ItemType:     config.DEDUCTION,
		IsDefault:    false,
		IsDeductible: true,
		IsTaxCost:    true,
		Amount:       0,
		PayRollID:    m.ID,
		Title:        "Biaya Jabatan",
	}).Error; err != nil {
		return err
	}

	// LATE DEDUCTION
	if err := database.DB.Create(&PayRollItem{
		ItemType:     config.DEDUCTION,
		IsDefault:    false,
		IsDeductible: true,
		Amount:       0,
		PayRollID:    m.ID,
		Title:        "Pot. Terlambat",
	}).Error; err != nil {
		return err
	}

	// LATE NOT PRESENCE
	if err := database.DB.Create(&PayRollItem{
		ItemType:     config.DEDUCTION,
		IsDefault:    false,
		IsDeductible: true,
		Amount:       0,
		PayRollID:    m.ID,
		Title:        "Pot. Tidak Masuk",
	}).Error; err != nil {
		return err
	}

	// BPJS

	if err := m.BpjsCount(); err != nil {
		return err
	}

	// TAX INCOME
	if err := database.DB.Create(&PayRollItem{
		ItemType:     config.DEDUCTION,
		IsDefault:    false,
		IsDeductible: false,
		IsTax:        true,
		TaxAutoCount: true,
		Amount:       0,
		PayRollID:    m.ID,
		Title:        "PPH 21",
	}).Error; err != nil {
		return err
	}

	return nil
}

func (m *PayRoll) BpjsCount() error {
	if m.BpjsSetting == nil {
		return nil
	}
	totalSalaryAndAllowance := float64(0)
	m.GetItems()
	for _, item := range m.Items {
		if item.BpjsCounted && item.ItemType != config.DEDUCTION {
			totalSalaryAndAllowance += item.Amount
		}
	}

	if m.BpjsSetting.BpjsKesEnabled {
		employeerCost, employeeCost, _ := m.BpjsSetting.CalculateBPJSKes(totalSalaryAndAllowance)
		payrollItem := PayRollItem{
			ItemType:     config.DEDUCTION,
			IsDefault:    false,
			IsDeductible: false,
			Amount:       employeeCost,
			PayRollID:    m.ID,
			Title:        "BPJS Kesehatan",
			Tariff:       m.BpjsSetting.BpjsKesRateEmployee,
		}

		if err := database.DB.Create(&payrollItem).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Biaya BPJS Kesehatan",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeerCost,
			Tariff:        m.BpjsSetting.BpjsKesRateEmployer,
		}).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Pungutan BPJS Kesehatan",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeeCost,
			Tariff:        m.BpjsSetting.BpjsKesRateEmployee,
			DebtDeposit:   true,
		}).Error; err != nil {
			return err
		}
	}
	if m.BpjsSetting.BpjsTkJhtEnabled {
		employeerCost, employeeCost, _ := m.BpjsSetting.CalculateBPJSTkJht(totalSalaryAndAllowance)
		payrollItem := PayRollItem{
			ItemType:     config.DEDUCTION,
			IsDefault:    false,
			IsDeductible: false,
			Amount:       employeeCost,
			PayRollID:    m.ID,
			Title:        "BPJS Ketenagakerjaan JHT",
			Tariff:       m.BpjsSetting.BpjsTkJhtRateEmployee,
		}

		if err := database.DB.Create(&payrollItem).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Biaya BPJS Ketenagakerjaan JHT",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeerCost,
			Tariff:        m.BpjsSetting.BpjsTkJhtRateEmployer,
			BpjsTkJht:     true,
		}).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Pungutan BPJS Ketenagakerjaan JHT",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeeCost,
			Tariff:        m.BpjsSetting.BpjsTkJhtRateEmployee,
			BpjsTkJht:     true,
			DebtDeposit:   true,
		}).Error; err != nil {
			return err
		}
	}
	if m.BpjsSetting.BpjsTkJpEnabled {
		employeerCost, employeeCost, _ := m.BpjsSetting.CalculateBPJSTkJp(totalSalaryAndAllowance)
		payrollItem := PayRollItem{
			ItemType:     config.DEDUCTION,
			IsDefault:    false,
			IsDeductible: false,
			Amount:       employeeCost,
			PayRollID:    m.ID,
			Title:        "BPJS Ketenagakerjaan JP",
			Tariff:       m.BpjsSetting.BpjsTkJpRateEmployee,
		}

		if err := database.DB.Create(&payrollItem).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Biaya BPJS Ketenagakerjaan JP",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeerCost,
			Tariff:        m.BpjsSetting.BpjsTkJpRateEmployer,
			BpjsTkJp:      true,
		}).Error; err != nil {
			return err
		}
		if err := database.DB.Create(&PayRollCost{
			Description:   "Pungutan BPJS Ketenagakerjaan JP",
			PayRollID:     m.ID,
			PayRollItemID: payrollItem.ID,
			Amount:        employeeCost,
			Tariff:        m.BpjsSetting.BpjsTkJpRateEmployee,
			BpjsTkJp:      true,
			DebtDeposit:   true,
		}).Error; err != nil {
			return err
		}
	}
	if m.BpjsSetting.BpjsTkJkmEnabled {
		employeeCost := m.BpjsSetting.CalculateBPJSTkJkm(totalSalaryAndAllowance)
		payrollItem := PayRollItem{
			ItemType:     config.DEDUCTION,
			IsDefault:    false,
			IsDeductible: false,
			Amount:       employeeCost,
			PayRollID:    m.ID,
			Title:        "BPJS Ketenagakerjaan JKM",
			Tariff:       m.BpjsSetting.BpjsTkJkmEmployee,
		}

		if err := database.DB.Create(&payrollItem).Error; err != nil {
			return err
		}
	}
	if m.BpjsSetting.BpjsTkJkkEnabled {
		employeeCost, tariff := m.BpjsSetting.CalculateBPJSTkJkk(totalSalaryAndAllowance, m.Employee.WorkSafetyRisks)
		payrollItem := PayRollItem{
			ItemType:     config.DEDUCTION,
			IsDefault:    false,
			IsDeductible: false,
			Amount:       employeeCost,
			PayRollID:    m.ID,
			Title:        "BPJS Ketenagakerjaan JKK",
			Tariff:       tariff,
		}

		if err := database.DB.Create(&payrollItem).Error; err != nil {
			return err
		}
	}

	return nil
}
func (m *PayRoll) GetDeductible() (float64, float64, float64, float64) {
	var totalIncome, totalReimbursement, totalDeductible, totalNonDeductible float64 = 0, 0, 0, 0
	for _, item := range m.Items {
		if item.IsTax {
			continue
		}
		if item.ItemType == config.DEDUCTION {
			if item.IsDeductible {
				totalDeductible += item.Amount
			} else {
				totalNonDeductible += item.Amount
			}
		} else {
			if item.ItemType == config.REIMBURSEMENT {
				totalReimbursement += item.Amount
			} else {
				totalIncome += item.Amount
			}
		}
	}

	return totalIncome, totalReimbursement, totalDeductible, totalNonDeductible
}

func (m *PayRoll) GetNonDeductibleItems() (items []PayRollItem) {
	for _, item := range m.Items {
		if item.IsTax {
			continue
		}
		if item.ItemType == config.DEDUCTION {
			if !item.IsDeductible {
				items = append(items, item)
			} else {
				continue
			}
		} else {
			continue
		}
	}
	return items
}

func (m *PayRoll) GetReimbursementItems() (items []PayRollItem) {
	for _, item := range m.Items {
		if item.ItemType == config.REIMBURSEMENT {
			items = append(items, item)
		}
	}
	return items
}

func (m *PayRoll) CountDeductible() {
	totalIncome, totalReimbursement, totalDeductible, totalNonDeductible := m.GetDeductible()
	// util.LogJson(map[string]interface{}{
	// 	"totalDeductible":    totalDeductible,
	// 	"totalNonDeductible": totalNonDeductible,
	// 	"totalReimbursement": totalReimbursement,
	// })
	m.TotalReimbursement = totalReimbursement
	m.TotalIncome = totalIncome + m.TaxAllowance
	m.TotalDeduction = totalDeductible + totalNonDeductible
	m.NetIncomeBeforeTaxCost = m.TotalIncome - totalDeductible
}

func (m *PayRoll) GetTaxCost(nonTaxable float64) float64 {
	taxCost := float64(0)
	if nonTaxable > 0 {
		taxCost = (m.NetIncomeBeforeTaxCost) * 5 / 100
		if taxCost > 500000 {
			taxCost = 500000
		}
	}
	return taxCost
}

// func (m *PayRoll) DeleteTaxAllowance() error {
// 	taxAllowanceItem := PayRollItem{}
// 	return database.DB.Delete(&taxAllowanceItem, "pay_roll_id = ? and is_tax_allowance = true", m.ID).Error
// }

func (m *PayRoll) EffectiveRateAverageTariff(category string, grossSalary float64) {
	ac := accounting.Accounting{Symbol: "", Precision: 4}
	taxTariff := float64(0)
	switch category {
	case "A":
		taxTariff = m.Employee.EffectiveRateAverageCategoryA(grossSalary)
	case "B":
		taxTariff = m.Employee.EffectiveRateAverageCategoryB(grossSalary)
	case "C":
		taxTariff = m.Employee.EffectiveRateAverageCategoryC(grossSalary)
	default:
		taxTariff = 0
	}
	taxAmount := grossSalary * taxTariff
	fmt.Printf("GROSS SALARY %s * TAXTARIFF %s = TAXAMOUNT %s \n", ac.FormatMoney(grossSalary), ac.FormatMoney(taxTariff), ac.FormatMoney(taxAmount))
	m.TotalTax = taxAmount
	m.TaxTariff = taxTariff
}

func (m *PayRoll) RefreshTax() {
	fmt.Println("RESET AMOUNT")
	database.DB.Model(&PayRoll{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"total_income":               0,
		"total_reimbursement":        0,
		"total_deduction":            0,
		"total_tax":                  0,
		"net_income_before_tax_cost": 0,
		"net_income":                 0,
		"take_home_pay":              0,
		"total_payable":              0,
		"tax_cost":                   0,
	})
	database.DB.Model(&PayRollItem{}).Where("pay_roll_id = ? AND is_tax_cost = true", m.ID).Update("amount", 0)
	database.DB.Model(&PayRollItem{}).Where("pay_roll_id = ? AND is_tax = true and tax_auto_count = true", m.ID).Update("amount", 0)
	if !m.IsGrossUp {
		database.DB.Model(&m).Where("id = ?", m.ID).Update("tax_allowance", 0)
		m.TaxAllowance = 0
	}
	m.TotalIncome = 0
	m.TotalReimbursement = 0
	m.TotalDeduction = 0
	m.TotalTax = 0
	m.NetIncomeBeforeTaxCost = 0
	m.NetIncome = 0
	m.TakeHomePay = 0
	m.TotalPayable = 0
	m.TaxCost = 0
}

func (m *PayRoll) GetItems() {
	items := []PayRollItem{}
	database.DB.Order("created_at asc").Find(&items, "pay_roll_id = ?", m.ID)
	m.Items = items
}
func (m *PayRoll) GetPayRollCost() {
	items := []PayRollCost{}
	database.DB.Order("created_at asc").Find(&items, "pay_roll_id = ?", m.ID)
	m.Costs = items
}
func (m *PayRoll) GetTotalPayRollCost() float64 {
	items := []PayRollCost{}
	database.DB.Order("created_at asc").Where("debt_deposit", 0).Find(&items, "pay_roll_id = ?", m.ID)
	// m.Items = items
	totalCost := float64(0)
	for _, v := range items {
		totalCost += v.Amount
	}

	return totalCost
}

func (m *PayRoll) GetTransactions() {
	transactions := []Transaction{}
	database.DB.Preload("AccountSource").Preload("AccountDestination").Joins("LEFT JOIN accounts as acc_source ON acc_source.id = transactions.account_source_id").Joins("LEFT JOIN accounts as acc_dest ON acc_dest.id = transactions.account_destination_id").Select("transactions.*, acc_source.is_tax is_source_tax, acc_source.name account_source_name, acc_dest.name account_destination_name, acc_dest.is_tax is_destination_tax").Find(&transactions, "transactions.pay_roll_id = ?", m.ID)
	m.Transactions = transactions
}

func (m *PayRoll) GetPayableTransactions() {
	transactions := []Transaction{}
	database.DB.Preload("AccountSource").Preload("AccountDestination").Joins("LEFT JOIN accounts as acc_source ON acc_source.id = transactions.account_source_id").Joins("LEFT JOIN accounts as acc_dest ON acc_dest.id = transactions.account_destination_id").Select("transactions.*, acc_source.is_tax is_source_tax, acc_source.name account_source_name, acc_dest.name account_destination_name, acc_dest.is_tax is_destination_tax").Where("(transactions.is_account_payable = true and transactions.tax_payment_id is  null and acc_dest.is_tax = false) or (transactions.pay_roll_payable_id is not null )").Find(&transactions, "transactions.pay_roll_id = ? and transactions.reimbursement_id is null", m.ID)
	m.PayableTransactions = transactions
}

func (m *PayRoll) RegularTaxTariff(taxAmount float64, taxable float64) {
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	// LEVEL 1
	if taxable > 0 {
		amountForTax := taxable - 60000000
		if amountForTax < 0 {
			amountForTax = taxable
			taxable = 0
		} else if amountForTax > 60000000 {
			taxable = amountForTax
			amountForTax = 60000000
		} else {
			taxable = amountForTax
			amountForTax = 60000000
		}
		taxValue := amountForTax * 5 / 100
		taxAmount += taxValue
		util.LogJson(map[string]interface{}{
			"msg":                "Level 1 => 5% add tax",
			"amountForTax":       ac.FormatMoney(amountForTax),
			"taxValue":           ac.FormatMoney(taxValue),
			"taxValue per month": ac.FormatMoney(taxValue / 12),
			"taxable":            ac.FormatMoney(taxable),
			"currentTaxAmount":   ac.FormatMoney(taxAmount),
		})
	}
	// LEVEL 2
	if taxable > 0 {
		amountForTax := taxable - 250000000
		if amountForTax < 0 {
			amountForTax = taxable
			taxable = 0
		} else if amountForTax > 250000000 {
			taxable = amountForTax
			amountForTax = 250000000
		} else {
			taxable = amountForTax
			amountForTax = 250000000
		}
		taxValue := amountForTax * 15 / 100
		taxAmount += taxValue
		util.LogJson(map[string]interface{}{
			"msg":                "Level 2 => 15% add tax",
			"amountForTax":       ac.FormatMoney(amountForTax),
			"taxValue":           ac.FormatMoney(taxValue),
			"taxValue per month": ac.FormatMoney(taxValue / 12),
			"taxable":            ac.FormatMoney(taxable),
			"currentTaxAmount":   ac.FormatMoney(taxAmount),
		})
	}
	// LEVEL 3
	if taxable > 0 {
		amountForTax := taxable - 500000000
		if amountForTax < 0 {
			amountForTax = taxable
			taxable = 0
		} else if amountForTax > 500000000 {
			taxable = amountForTax
			amountForTax = 500000000
		} else {
			taxable = amountForTax
			amountForTax = 500000000
		}
		taxValue := amountForTax * 25 / 100
		taxAmount += taxValue
		util.LogJson(map[string]interface{}{
			"msg":                "Level 3 => 25% add tax",
			"amountForTax":       ac.FormatMoney(amountForTax),
			"taxValue":           ac.FormatMoney(taxValue),
			"taxValue per month": ac.FormatMoney(taxValue / 12),
			"taxable":            ac.FormatMoney(taxable),
			"currentTaxAmount":   ac.FormatMoney(taxAmount),
		})
	}
	// LEVEL 4
	if taxable > 0 {
		amountForTax := taxable - 5000000000
		if amountForTax < 0 {
			amountForTax = taxable
			taxable = 0
		} else if amountForTax > 5000000000 {
			taxable = amountForTax
			amountForTax = 5000000000
		} else {
			taxable = amountForTax
			amountForTax = 5000000000
		}
		taxValue := amountForTax * 30 / 100
		taxAmount += taxValue
		util.LogJson(map[string]interface{}{
			"msg":                "Level 4 => 30% add tax",
			"amountForTax":       ac.FormatMoney(amountForTax),
			"taxValue":           ac.FormatMoney(taxValue),
			"taxValue per month": ac.FormatMoney(taxValue / 12),
			"taxable":            ac.FormatMoney(taxable),
			"currentTaxAmount":   ac.FormatMoney(taxAmount),
		})
	}
	// LEVEL 5
	if taxable > 0 {
		taxValue := taxable * 35 / 100
		taxAmount += taxValue
		util.LogJson(map[string]interface{}{
			"msg":                "Level 5 => 35% add tax",
			"taxValue":           ac.FormatMoney(taxValue),
			"taxValue per month": ac.FormatMoney(taxValue / 12),
			"taxable":            ac.FormatMoney(taxable),
			"currentTaxAmount":   ac.FormatMoney(taxAmount),
		})
	}
	util.LogJson(map[string]interface{}{
		"taxAmount":           ac.FormatMoney(taxAmount),
		"taxAmount per month": ac.FormatMoney(taxAmount / 12),
	})
	m.TotalTax = taxAmount / 12
}

func (m *PayRoll) CountTax() error {
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	m.RefreshTax()
	m.GetItems()
	m.GetEmployee()

	nonTaxable := m.Employee.GetNonTaxableIncomeLevelAmount()
	nonTaxableCategory := m.Employee.GetNonTaxableIncomeLevelCategory()

	countTaxRecord := int64(0)
	taxCost := float64(0)

	taxCostItem := PayRollItem{}
	// taxItem := PayRollItem{}
	database.DB.Find(&taxCostItem, "pay_roll_id = ? AND is_tax_cost = true", m.ID)
	// database.DB.Find(&taxItem, "pay_roll_id = ? AND is_tax_cost = true", m.ID)
	database.DB.Model(&PayRollItem{}).Where("pay_roll_id = ? AND is_tax = true and tax_auto_count = true", m.ID).Count(&countTaxRecord)
	m.CountDeductible()

	taxCost = m.GetTaxCost(nonTaxable)
	if m.IsEffectiveRateAverage {
		taxCost = 0
	}

	m.NetIncome = m.NetIncomeBeforeTaxCost - taxCost
	database.DB.Model(&taxCostItem).Where("id = ?", taxCostItem.ID).Update("amount", taxCost)

	database.DB.Model(&m).Update("tax_cost", taxCost)
	fmt.Println("TAX COST", ac.FormatMoney(taxCost))
	fmt.Println("NET INCOME AFTER REDUCE TAX COST", ac.FormatMoney(m.NetIncomeBeforeTaxCost))
	fmt.Println("TAX AUTO COUNT", countTaxRecord)
	// GET TAX COST
	// 1. GET NET INCOME
	m.GetItems()

	yearlyNetIncome := (m.NetIncome + m.GetTotalPayRollCost()) * 12

	taxable := yearlyNetIncome - nonTaxable

	util.LogJson(map[string]interface{}{
		"yearlyGrossIncome":    ac.FormatMoney(m.NetIncomeBeforeTaxCost * 12),
		"taxCost":              ac.FormatMoney(taxCost * 12),
		"yearlyNetIncome":      ac.FormatMoney(yearlyNetIncome),
		"taxable":              ac.FormatMoney(taxable),
		"nonTaxable":           ac.FormatMoney(nonTaxable),
		"taxAllowancePerMonth": ac.FormatMoney(m.TaxAllowance),
	})
	var taxAmount float64 = m.TotalTax
	if countTaxRecord == 0 {
		taxAmount = 0
		m.TotalTax = 0
		database.DB.Model(&m).Update("total_tax", 0)

	}

	if nonTaxable != 0 && countTaxRecord > 0 {
		if m.IsEffectiveRateAverage {
			fmt.Println("NON_TAXABLE_CATEGORY", nonTaxableCategory)
			m.EffectiveRateAverageTariff(nonTaxableCategory, m.NetIncomeBeforeTaxCost+m.GetTotalPayRollCost())
		} else {
			m.TaxTariff = 0
			m.RegularTaxTariff(taxAmount, taxable)
			database.DB.Model(&m).Update("tax_tariff", 0)

		}
	}

	m.TakeHomePay = m.TotalIncome - m.TotalDeduction - m.TotalTax - m.TaxCost
	database.DB.Model(&PayRollItem{}).Where("pay_roll_id = ? AND is_tax = true and tax_auto_count = true", m.ID).Update("amount", m.TotalTax)
	if err := database.DB.Model(&m).Updates(&m).Error; err != nil {
		return err
	}

	if m.TotalReimbursement == 0 {
		if err := database.DB.Model(&m).Update("total_reimbursement", 0).Error; err != nil {
			return err
		}
	}

	if m.IsGrossUp {
		fmt.Println("m.TaxAllowance", m.TaxAllowance)
		fmt.Println("m.TotalTax", m.TotalTax)
		if m.TotalTax != m.TaxAllowance {
			m.TaxAllowance = m.TotalTax
			if err := database.DB.Model(&m).Updates(&m).Error; err != nil {
				return err
			}
			fmt.Println("m.TaxAllowance UPDATED", m.TaxAllowance)
			return m.CountTax()
		}
	}

	m.TaxSummary = CountTaxSummary{
		JobExpenseMonthly:               taxCost,
		JobExpenseYearly:                taxCost * 12,
		PtkpYearly:                      nonTaxable,
		GrossIncomeMonthly:              m.NetIncomeBeforeTaxCost,
		GrossIncomeYearly:               m.NetIncomeBeforeTaxCost * 12,
		PkpMonthly:                      (m.NetIncome*12 - nonTaxable) / 12,
		PkpYearly:                       m.NetIncome*12 - nonTaxable,
		TaxYearlyBasedOnProgressiveRate: m.TotalTax * 12,
		TaxYearly:                       m.TotalTax * 12,
		TaxMonthly:                      m.TotalTax,
		NetIncomeMonthly:                m.NetIncome,
		NetIncomeYearly:                 m.NetIncome,
		CutoffPensiunMonthly:            0,
		CutoffPensiunYearly:             0,
		CutoffMonthly:                   0,
		CutoffYearly:                    0,
		Ter:                             0,
	}

	return nil

}

type CountTaxSummary struct {
	JobExpenseMonthly               float64 `json:"job_expense_monthly"`
	JobExpenseYearly                float64 `json:"job_expense_yearly"`
	PtkpYearly                      float64 `json:"ptkp_yearly"`
	GrossIncomeMonthly              float64 `json:"gross_income_monthly"`
	GrossIncomeYearly               float64 `json:"gross_income_yearly"`
	PkpMonthly                      float64 `json:"pkp_monthly"`
	PkpYearly                       float64 `json:"pkp_yearly"`
	TaxYearlyBasedOnProgressiveRate float64 `json:"tax_yearly_based_on_progressive_rate"`
	TaxYearly                       float64 `json:"tax_yearly"`
	TaxMonthly                      float64 `json:"tax_monthly"`
	NetIncomeMonthly                float64 `json:"net_income_monthly"`
	NetIncomeYearly                 float64 `json:"net_income_yearly"`
	CutoffPensiunMonthly            float64 `json:"cutoff_pensiun_monthly"`
	CutoffPensiunYearly             float64 `json:"cutoff_pensiun_yearly"`
	CutoffMonthly                   float64 `json:"cutoff_monthly"`
	CutoffYearly                    float64 `json:"cutoff_yearly"`
	Ter                             float64 `json:"ter"`
}
