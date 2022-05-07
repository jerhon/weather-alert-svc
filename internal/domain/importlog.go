package domain

import "time"

type ImportLog struct {
	Type           string
	ImportedTime   time.Time
	LastModified   string
	ImportedAlerts int
}
