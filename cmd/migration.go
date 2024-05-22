package cmd

import (
	"avolta/config"
	"avolta/database"
	"avolta/model"
	"avolta/object/constants"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-faker/faker/v4"
)

func Migrate() {
	fmt.Println("START TO MIGRATE")
	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Employee{})
	database.DB.AutoMigrate(&model.PayRoll{})
	database.DB.AutoMigrate(&model.Attendance{})
	database.DB.AutoMigrate(&model.Incentive{})
	database.DB.AutoMigrate(&model.Role{})
	database.DB.AutoMigrate(&model.Permission{})
	database.DB.AutoMigrate(&model.Image{})
	database.DB.AutoMigrate(&model.PayRoll{})
	database.DB.AutoMigrate(&model.PayRollItem{})
	database.DB.AutoMigrate(&model.Permission{})
	database.DB.AutoMigrate(&model.Role{})
	database.DB.AutoMigrate(&model.Transaction{})
	database.DB.AutoMigrate(&model.Account{})
	database.DB.AutoMigrate(&model.Organization{})
	database.DB.AutoMigrate(&model.JobTitle{})
	database.DB.AutoMigrate(&model.Schedule{})
	database.DB.AutoMigrate(&model.AttendanceBulkImport{})
	database.DB.AutoMigrate(&model.AttendanceImport{})
	database.DB.AutoMigrate(&model.AttendanceImportItem{})
	database.DB.AutoMigrate(&model.LeaveCategory{})
	database.DB.AutoMigrate(&model.Leave{})
	database.DB.AutoMigrate(&model.Company{})
	database.DB.AutoMigrate(&model.Setting{})

	fmt.Println("FINISHED  MIGRATE")
}

func TestCreateUser(args []string) {
	database.DB.Create(&model.User{
		Email:    args[1],
		Password: args[2],
		IsAdmin:  true,
	})
}

func SampleAttendance(args []string) {
	if len(args) == 0 {
		fmt.Println("please set number of attendance")
		os.Exit(0)
	}

	max, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	var employees []model.Employee
	database.DB.Limit(20).Order("rand()").Find(&employees)

	for i := 0; i < max; i++ {
		rand.New(rand.NewSource(time.Now().Unix()))
		index := rand.Intn(len(employees))
		min := 00
		max := 60
		hourMin := 1
		hourMax := 9
		dateMin := 1
		dateMax := 15
		g := fmt.Sprintf("%02d", rand.Intn(dateMax-dateMin)+dateMin)

		clockin, _ := time.Parse("2006-01-02 15:04:05", "2024-05-"+g+" "+faker.TimeString())
		clockout := clockin.Add(time.Hour*time.Duration(rand.Intn(hourMax-hourMin)+hourMin) + time.Minute*time.Duration(rand.Intn(max-min)+min))
		database.DB.Create(&model.Attendance{
			EmployeeID:    &employees[index].ID,
			ClockIn:       clockin,
			ClockOut:      &clockout,
			ClockInNotes:  faker.Word(),
			ClockOutNotes: faker.Word(),
		})
	}
}
func SampleJobTitle() {
	positions := []string{"Staff", "Manager", "Leader", "Supervisor", "HRD"}
	for _, v := range positions {
		database.DB.Create(&model.JobTitle{
			Name: v,
		})
	}
}
func SampleEmployee(args []string) {
	if len(args) == 0 {
		fmt.Println("please set number of employee")
		os.Exit(0)
	}

	max, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	jobTitles := []model.JobTitle{}
	database.DB.Find(&jobTitles)
	for i := 0; i < max; i++ {

		rand.New(rand.NewSource(time.Now().Unix()))
		genders := []string{"f", "m"}
		randomIndex := rand.Intn(len(genders))

		gender := genders[randomIndex]
		randomIndex2 := rand.Intn(len(jobTitles))
		JobTitleID := jobTitles[randomIndex2].ID
		var startedWork, birthDate time.Time
		if startedWork, err = time.Parse("2006-01-02", faker.Date()); err != nil {
			startedWork = time.Time{}
		}

		if birthDate, err = time.Parse("2006-01-02", faker.Date()); err != nil {
			birthDate = time.Time{}
		}

		database.DB.Create(&model.Employee{
			Email:       faker.Email(),
			Gender:      gender,
			FirstName:   faker.FirstName(),
			MiddleName:  faker.Word(),
			LastName:    faker.LastName(),
			Phone:       faker.Phonenumber(),
			JobTitleID:  model.NullStringConv(JobTitleID),
			DateOfBirth: model.NullTimeConv(birthDate),
			StartedWork: model.NullTimeConv(startedWork),
		})
	}
}
func AssignSuperadmin(args []string) {
	user := model.User{}
	role := model.Role{}
	if err := database.DB.Find(&user, "email = ?", args[1]).Error; err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if err := database.DB.Find(&role, "is_super_admin = 1").Error; err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if err := database.DB.Model(&user).Update("role_id", role.ID).Error; err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func GenAccounts() {

	// CURRENT ASSET
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.CASH_BANK,
		CashflowGroup:    config.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         config.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Kecil",
		Type:             config.TYPE_ASSET,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.CASH_BANK,
		CashflowGroup:    config.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         config.CATEGORY_CURRENT_ASSET,
		Name:             "BANK",
		Type:             config.TYPE_ASSET,
	})

	// EQUITY
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.EQUITY_CAPITAL,
		CashflowGroup:    config.CASHFLOW_GROUP_FINANCING,
		Category:         config.CATEGORY_EQUITY,
		Name:             "Modal Awal",
		Type:             config.TYPE_EQUITY,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.EQUITY_CAPITAL,
		CashflowGroup:    config.CASHFLOW_GROUP_FINANCING,
		Category:         config.CATEGORY_EQUITY,
		Name:             "Ekuitas Saldo Awal",
		Type:             config.TYPE_EQUITY,
		IsDeletable:      true,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.EQUITY_CAPITAL,
		CashflowGroup:    config.CASHFLOW_GROUP_FINANCING,
		Category:         config.CATEGORY_EQUITY,
		Name:             "Saldo Awal Aset Tetap",
		Type:             config.TYPE_EQUITY,
		IsDeletable:      true,
	})
	// INCOME
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.ACCEPTANCE_FROM_CUSTOMERS,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_SALES,
		Name:             "Penjualan",
		Type:             config.TYPE_INCOME,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Retur Penjualan",
		Type:             config.TYPE_EXPENSE,
	})

	// EXPENSE
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Pembayaran Gaji",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Peralatan Kantor",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Pembayaran Listrik",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Rekening Telepon & Pulsa",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})

	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Internet",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})

	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_OPERATING,
		Name:             "Bonus Pegawai",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_EXPENSE,
		Name:             "Beban Lainnya",
		IsDeletable:      true,
		Type:             config.TYPE_EXPENSE,
	})

	database.DB.Create(&model.Account{
		CashflowSubGroup: config.ACQUISITION_SALE_OF_ASSETS,
		CashflowGroup:    config.CASHFLOW_GROUP_INVESTING,
		Category:         config.CATEGORY_EXPENSE,
		Name:             "Penyusutan Tanah - Bangunan",
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.ACQUISITION_SALE_OF_ASSETS,
		CashflowGroup:    config.CASHFLOW_GROUP_INVESTING,
		Category:         config.CATEGORY_EXPENSE,
		Name:             "Penyusutan Kendaraan",
		Type:             config.TYPE_EXPENSE,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.ACQUISITION_SALE_OF_ASSETS,
		CashflowGroup:    config.CASHFLOW_GROUP_INVESTING,
		Category:         config.CATEGORY_EXPENSE,
		Name:             "Penyusutan Lainnya",
		Type:             config.TYPE_EXPENSE,
	})

	// RECEIVABLE
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.ACCEPTANCE_FROM_CUSTOMERS,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_RECEIVABLE,
		Name:             "Piutang Usaha",
		Type:             config.TYPE_RECEIVABLE,
	})

	// LIABILITY
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.PAYMENT_TO_VENDORS,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Usaha",
		Type:             config.TYPE_LIABILITY,
	})
	database.DB.Create(&model.Account{
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Gaji / Reimbursement",
		Type:             config.TYPE_LIABILITY,
	})
	// COST
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_COST_OF_REVENUE,
		Name:             "Beban Pokok Pendapatan",
		Type:             config.TYPE_COST,
	})
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_PURCHASE_RETURNS,
		Name:             "Retur Pembelian",
		Type:             config.TYPE_COST,
	})

	database.DB.Create(&model.Account{
		CashflowSubGroup: config.OPERATIONAL_EXPENSES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_PRODUCTION_COST,
		Name:             "Biaya Produksi",
		Type:             config.TYPE_COST,
	})

	// Hutang Pajak - PPh 21
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.RETURNS_PAYMENT_OF_TAXES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Pajak - PPh 21",
		IsTax:            true,
		Type:             config.TYPE_LIABILITY,
	})

	// Hutang Pajak - PPh 22
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.RETURNS_PAYMENT_OF_TAXES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Pajak - PPh 22",
		IsTax:            true,
		Type:             config.TYPE_LIABILITY,
	})
	// Hutang Pajak - PPh 23
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.RETURNS_PAYMENT_OF_TAXES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Pajak - PPh 23",
		IsTax:            true,
		Type:             config.TYPE_LIABILITY,
	})
	// Hutang Pajak - PPh 29
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.RETURNS_PAYMENT_OF_TAXES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Pajak - PPh 29",
		IsTax:            true,
		Type:             config.TYPE_LIABILITY,
	})

	// Hutang Pajak Lainnya
	database.DB.Create(&model.Account{
		CashflowSubGroup: config.RETURNS_PAYMENT_OF_TAXES,
		CashflowGroup:    config.CASHFLOW_GROUP_OPERATING,
		Category:         config.CATEGORY_DEBT,
		Name:             "Hutang Pajak Lainnya",
		IsTax:            true,
		Type:             config.TYPE_LIABILITY,
	})

}

func GenPermissions() {
	permissions := constants.DefaultPermission("")
	for _, v := range permissions {
		if err := database.DB.Create(&v).Error; err != nil {
			fmt.Println("ERROR CREATE PERMISSION ", v.Name, err)
		}
	}
	GenSuperAdmin()
}
func GenLeaveCategories() {
	cats := []string{
		"Izin Sakit",
		"Sakit dengan Surat Dokter",
		"Dinas Luar Kota",
		"Cuti Menikah",
		"Cuti Menikahkan Anak",
		"Cuti Khitanan Anak",
		"Cuti Baptis Anak",
		"Cuti Istri Melahirkan atau Keguguran",
		"Cuti Keluarga Meninggal",
		"Cuti Anggota Keluarga Dalam Satu Rumah Meninggal",
		"Cuti Ibadah Haji",
		"Izin Lainnya",
	}

	for _, v := range cats {
		database.DB.Create(&model.LeaveCategory{
			Name: v,
		})
	}
}
func GenSuperAdmin() {
	database.DB.Create(&model.Role{
		Name:         "SUPERADMIN",
		Description:  "Yes i'am superman",
		IsSuperAdmin: true,
	})
}
