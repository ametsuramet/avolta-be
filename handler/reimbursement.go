package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReimbursementGetAllHandler(c *gin.Context) {
	var data []model.Reimbursement
	preloads := []string{"Employee"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	employeeId, ok := c.GetQuery("employee_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"employee_id": employeeId,
		})

	}
	status, ok := c.GetQuery("status")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"status": status,
		})

	}

	paginator.Paginate(&data)
	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Reimbursement Retrived", dataRecords.Records, dataRecords)
}

func ReimbursementGetOneHandler(c *gin.Context) {
	var data model.Reimbursement

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Items").Preload("Transactions", func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at asc")
		return db
	}).Preload("Transactions.AccountDestination").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	total := float64(0)
	totalPaid := float64(0)
	for _, v := range data.Items {
		total += v.Amount
	}

	if err := database.DB.Model(&data).Update("total", total).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	for _, v := range data.Transactions {
		if v.IsReimbursementPayment {
			totalPaid += v.Debit
		}
	}
	if err := database.DB.Model(&data).Update("balance", total-totalPaid).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Reimbursement Retrived", data, nil)
}

func ReimbursementCreateHandler(c *gin.Context) {
	var data model.Reimbursement

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Reimbursement Created", gin.H{"last_id": data.ID}, nil)
}

func ReimbursementUpdateHandler(c *gin.Context) {
	var input, data model.Reimbursement
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Reimbursement Updated", nil, nil)
}

func ReimbursementApprovalHandler(c *gin.Context) {
	var data model.Reimbursement
	id := c.Params.ByName("id")
	approvalType := c.Params.ByName("type")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := approvalReimbursement(c, approvalType); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data Reimbursement Updated", nil, nil)
}
func ReimbursemenPaymentHandler(c *gin.Context) {
	var data model.Reimbursement
	id := c.Params.ByName("id")
	var input = struct {
		Remarks   string  `json:"remarks"`
		Amount    float64 `json:"amount" `
		AccountID string  `json:"account_id" binding:"required"`
		Files     string  `json:"files"`
	}{}

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now()
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := c.ShouldBindJSON(&input); err != nil {
			return err
		}
		var setting model.Setting
		if err := database.DB.First(&setting).Error; err != nil {
			return err
		}

		getUser, _ := c.Get("user")
		user := getUser.(model.User)
		balance := data.Balance - input.Amount

		if balance == 0 {
			data.Status = "PAID"
		}

		files := []string{}
		newFiles := []string{}

		if data.Attachment != "" {
			json.Unmarshal([]byte(data.Attachment), &files)

		}
		if input.Files != "" {
			json.Unmarshal([]byte(input.Files), &newFiles)

		}
		files = append(files, newFiles...)
		b, _ := json.Marshal(&files)
		data.Attachment = string(b)

		data.Remarks += "- **[" + user.FullName + "]** \n*" + time.Now().Format("02/01/2006 15:04") + "*\n\n **" + "PAYMENT" + "** " + input.Remarks + "\n\n"
		if err := tx.Model(&data).Updates(&data).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Transaction{
			Description:          "Pembayaran Hutang " + data.Name + "",
			Debit:                input.Amount,
			AccountDestinationID: setting.ReimbursementPayableAccountID,
			IsAccountPayable:     true,
			Date:                 now,
			ReimbursementID:      &data.ID,
			EmployeeID:           data.EmployeeID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Transaction{
			Description:            "Pembayaran " + data.Name + "",
			Debit:                  input.Amount,
			AccountDestinationID:   &input.AccountID,
			Date:                   now,
			ReimbursementID:        &data.ID,
			EmployeeID:             data.EmployeeID,
			IsReimbursementPayment: true,
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&data).Update("balance", balance).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data Reimbursement Updated", nil, nil)
}

func ReimbursementDeleteHandler(c *gin.Context) {
	var input, data model.Reimbursement
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Reimbursement Deleted", nil, nil)
}

func approvalReimbursement(c *gin.Context, status string) error {
	var data model.Reimbursement
	var input = struct {
		Remarks string `json:"remarks"`
	}{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		return err
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := database.DB.Transaction(func(tx *gorm.DB) error {

		getUser, _ := c.Get("user")
		user := getUser.(model.User)

		data.Status = status
		data.Remarks += "- **[" + user.FullName + "]** \n*" + time.Now().Format("02/01/2006 15:04") + "*\n\n **" + status + "** " + input.Remarks + "\n\n"
		if err := tx.Model(&data).Updates(&data).Error; err != nil {
			return err
		}
		var setting model.Setting
		if err := database.DB.First(&setting).Error; err != nil {
			return err
		}

		if status == "APPROVED" {
			if err := tx.Create(&model.Transaction{

				Description:          data.Name,
				Debit:                data.Total,
				AccountDestinationID: setting.ReimbursementExpenseAccountID,
				IsExpense:            true,
				Date:                 now,
				ReimbursementID:      &data.ID,
				EmployeeID:           data.EmployeeID,
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&model.Transaction{

				Description:          "Hutang " + data.Name + "",
				Credit:               data.Total,
				AccountDestinationID: setting.ReimbursementPayableAccountID,
				IsAccountPayable:     true,
				Date:                 now,
				ReimbursementID:      &data.ID,
				EmployeeID:           data.EmployeeID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
