package cmd

func Init(args []string) {

	switch args[0] {
	case "migrate":
		Migrate()
	case "user":
		TestCreateUser(args)
	case "gen-account":
		GenAccounts()
	case "gen-permission":
		GenPermissions()
	case "assign-superadmin":
		AssignSuperadmin(args)
	case "sample-employee":
		SampleEmployee(args)
	case "sample-attendance":
		SampleAttendance(args)
	case "gen":
		GenerateFeature(args)

	}
}
