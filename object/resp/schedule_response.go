package resp

type ScheduleReponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ScheduleType string            `json:"schedule_type" `
	WeekDay      string            `json:"week_day" `
	StartDate    string            `json:"start_date" `
	EndDate      string            `json:"end_date" `
	StartTime    string            `json:"start_time" `
	EndTime      string            `json:"end_time" `
	Employees    []EmployeeReponse `json:"employees" `
	Sunday       bool              `json:"sunday"`
	Monday       bool              `json:"monday"`
	Tuesday      bool              `json:"tuesday"`
	Wednesday    bool              `json:"wednesday"`
	Thursday     bool              `json:"thursday"`
	Friday       bool              `json:"friday"`
	Saturday     bool              `json:"saturday"`
}
