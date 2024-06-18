package cmd

func Init(args []string) {

	switch args[0] {
	case "migrate":
		Migrate()
	case "user":
		TestCreateUser(args)
	// case "gen-account":
	// 	GenAccounts(args)
	case "gen-permission":
		GenPermissions()
	case "gen-bank":
		GenBanks()
	case "gen-leave-category":
		GenLeaveCategories()
	case "gen-product-category":
		GenProductCategories()
	case "assign-superadmin":
		AssignSuperadmin(args)
	case "sample-job-title":
		SampleJobTitle()
	case "sample-employee":
		SampleEmployee(args)
	case "sample-attendance":
		SampleAttendance(args)
	case "gen":
		GenerateFeature(args)

	}
}
