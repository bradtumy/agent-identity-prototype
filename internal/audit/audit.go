package audit

import "log"

// LogAction logs an action for auditing purposes.
func LogAction(action string, subject string, success bool) {
	log.Printf("AUDIT action=%s subject=%s success=%t", action, subject, success)
}
