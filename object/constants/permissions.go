package constants

import (
	"avolta/model"
	"strings"
)

func DefaultPermission(search string) []model.Permission {
	var permissions []model.Permission
	var cruds = []string{"create", "read", "update", "delete"}

	var features = []map[string]interface{}{
		{"name": "account", "is_default": true, "is_active": true},
		{"name": "schedule", "is_default": true, "is_active": true},
		{"name": "job_title", "is_default": true, "is_active": true},
		{"name": "attendance", "is_default": true, "is_active": true},
		{"name": "reimbursement", "is_default": true, "is_active": true},
		{"name": "employee", "is_default": true, "is_active": true},
		{"name": "transaction", "is_default": true, "is_active": true},
		{"name": "incentive", "is_default": true, "is_active": true},
		{"name": "leave", "is_default": true, "is_active": true},
		{"name": "leave_category", "is_default": true, "is_active": true},
		{"name": "overtime", "is_default": true, "is_active": true},
		{"name": "pay_roll", "is_default": true, "is_active": true},
		{"name": "role", "is_default": true, "is_active": true},
		{"name": "user", "is_default": true, "is_active": true},
		{"name": "organization", "is_default": true, "is_active": true},
		{"name": "product", "is_default": true, "is_active": true},
		{"name": "product_category", "is_default": true, "is_active": true},
		{"name": "sale", "is_default": true, "is_active": true},
		{"name": "sale_receipt", "is_default": true, "is_active": true},
		{"name": "shop", "is_default": true, "is_active": true},
		{"name": "incentive_shop", "is_default": true, "is_active": true},
		{"name": "incentive_setting", "is_default": true, "is_active": true},
		{"name": "incentive_report", "is_default": true, "is_active": true},
	}
	var reports = []map[string]interface{}{
		{"name": "attendance", "is_default": true, "is_active": true},
		{"name": "pay_roll", "is_default": true, "is_active": true},
		{"name": "transaction", "is_default": true, "is_active": true},
		{"name": "cash_flow", "is_default": true, "is_active": true},
	}
	var menus = []map[string]interface{}{
		{"name": "dashboard", "is_default": true, "is_active": true},
		{"name": "employee", "is_default": true, "is_active": true},
		{"name": "attendance", "is_default": true, "is_active": true},
		{"name": "leave", "is_default": true, "is_active": true},
		{"name": "role", "is_default": true, "is_active": true},
		{"name": "pay_roll", "is_default": true, "is_active": true},
		{"name": "company", "is_default": true, "is_active": true},
		{"name": "report", "is_default": true, "is_active": true},
		{"name": "system", "is_default": true, "is_active": true},
	}
	var imports = []map[string]interface{}{
		{"name": "dashboard", "is_default": true, "is_active": true},
		{"name": "employee", "is_default": true, "is_active": true},
		{"name": "attendance", "is_default": true, "is_active": true},
		{"name": "product", "is_default": true, "is_active": true},
		{"name": "sale", "is_default": true, "is_active": true},
		{"name": "leave", "is_default": true, "is_active": true},
		{"name": "role", "is_default": true, "is_active": true},
		{"name": "pay_roll", "is_default": true, "is_active": true},
		{"name": "company", "is_default": true, "is_active": true},
		{"name": "incentive", "is_default": true, "is_active": true},
		{"name": "report", "is_default": true, "is_active": true},
	}

	for _, feature := range features {
		for _, crud := range cruds {
			var split = strings.Split(feature["name"].(string), "_")
			for i, _ := range split {
				split[i] = strings.Title(split[i])
			}
			var title = strings.Join(split, " ")
			var name = strings.Title(crud) + " " + title
			var key = crud + "_" + feature["name"].(string)
			if key == "create_company" {
				continue
			}
			permissions = append(permissions, model.Permission{
				Name:      name,
				Key:       key,
				IsDefault: feature["is_default"].(bool),
				IsActive:  feature["is_active"].(bool),
				Group:     feature["name"].(string),
			})
		}
	}

	for _, report := range reports {
		var split = strings.Split(report["name"].(string), "_")
		for i, _ := range split {
			split[i] = strings.Title(split[i])
		}
		var title = strings.Join(split, " ")
		permissions = append(permissions, model.Permission{
			Name:      "Report " + title,
			Key:       "report_" + report["name"].(string),
			IsDefault: report["is_default"].(bool),
			IsActive:  report["is_active"].(bool),
			Group:     "report",
		})
	}
	for _, menu := range menus {
		var split = strings.Split(menu["name"].(string), "_")
		for i, _ := range split {
			split[i] = strings.Title(split[i])
		}
		var title = strings.Join(split, " ")
		permissions = append(permissions, model.Permission{
			Name:      "Menu " + title,
			Key:       "menu_" + menu["name"].(string),
			IsDefault: menu["is_default"].(bool),
			IsActive:  menu["is_active"].(bool),
			Group:     "menu",
		})
	}
	for _, importDoc := range imports {
		var split = strings.Split(importDoc["name"].(string), "_")
		for i, _ := range split {
			split[i] = strings.Title(split[i])
		}
		var title = strings.Join(split, " ")
		permissions = append(permissions, model.Permission{
			Name:      "Import " + title,
			Key:       "import_" + importDoc["name"].(string),
			IsDefault: importDoc["is_default"].(bool),
			IsActive:  importDoc["is_active"].(bool),
			Group:     "import",
		})
	}

	permissions = append(permissions, model.Permission{
		Name:      "Import Attendance Approval",
		Key:       "import_attendance_approval",
		IsDefault: true,
		IsActive:  true,
		Group:     "import",
	})

	permissions = append(permissions, model.Permission{
		Name:      "Leave Approval",
		Key:       "approval_leave",
		IsDefault: true,
		IsActive:  true,
		Group:     "leave",
	})

	permissions = append(permissions, model.Permission{
		Name:      "Reimbursement Approval",
		Key:       "approval_reimbursement",
		IsDefault: true,
		IsActive:  true,
		Group:     "reimbursement",
	})
	permissions = append(permissions, model.Permission{
		Name:      "Reimbursement Payment",
		Key:       "payment_reimbursement",
		IsDefault: true,
		IsActive:  true,
		Group:     "reimbursement",
	})

	permissions = append(permissions, model.Permission{
		Name:      "Payment Pay Roll",
		Key:       "payment_pay_roll",
		IsDefault: true,
		IsActive:  true,
		Group:     "pay_roll",
	})

	permissions = append(permissions, model.Permission{
		Name:      "Payment Invoice",
		Key:       "payment_invoice",
		IsDefault: true,
		IsActive:  true,
		Group:     "invoice",
	})
	permissions = append(permissions, model.Permission{
		Name:      "User Log",
		Key:       "log_user",
		IsDefault: true,
		IsActive:  true,
		Group:     "miscellaneous",
	})

	permissions = append(permissions, model.Permission{
		Name:      "Incentive Approval",
		Key:       "approval_incentive",
		IsDefault: true,
		IsActive:  true,
		Group:     "incentive",
	})
	permissions = append(permissions, model.Permission{
		Name:      "Incentive Payment",
		Key:       "payment_incentive",
		IsDefault: true,
		IsActive:  true,
		Group:     "incentive",
	})

	if search != "" && len(search) > 1 {
		newPermissions := []model.Permission{}
		for _, v := range permissions {
			if strings.Contains(v.Key, search) {
				newPermissions = append(newPermissions, model.Permission{
					Name:      v.Name,
					Key:       v.Key,
					IsDefault: v.IsDefault,
					Group:     v.Group,
				})
			}
		}
		permissions = newPermissions
	}
	return permissions
}
