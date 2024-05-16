package resp

import "time"

type AttendanceReponse struct {
	ClockIn          time.Time  `json:"clock_in"`
	ClockOut         *time.Time `json:"clock_out"`
	ClockInNotes     string     `json:"clock_in_notes"`
	ClockOutNotes    string     `json:"clock_out_notes"`
	ClockInPicture   string     `json:"clock_in_picture"`
	ClockOutPicture  string     `json:"clock_out_picture"`
	ClockInLat       float64    `json:"clock_in_lat"`
	ClockInLng       float64    `json:"clock_in_lng"`
	ClockOutLat      float64    `json:"clock_out_lat"`
	ClockOutLng      float64    `json:"clock_out_lng"`
	EmployeeID       *string    `json:"employee_id"`
	EmployeeName     string     `json:"employee_name"`
	EmployeePosition string     `json:"employee_position"`
	EmployeePicture  *string    `json:"employee_picture"`
}
