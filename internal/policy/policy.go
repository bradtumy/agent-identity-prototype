package policy

import "errors"

var allowedActions = []string{"fetch_data", "transform", "notify"}

var rolePermissions = map[string][]string{
	"data-fetcher": {"fetch_data"},
	"transformer":  {"transform"},
	"notifier":     {"notify"},
}

func IsActionAllowedForRole(role, action string) bool {
	actions, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}

func ValidatePolicy(action, role string) error {
	if !contains(allowedActions, action) {
		return errors.New("action not allowed")
	}
	if !IsActionAllowedForRole(role, action) {
		return errors.New("role not permitted to perform action")
	}
	return nil
}

func contains(list []string, item string) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}
