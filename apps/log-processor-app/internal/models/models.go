package models

import (
	"time"
)

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

type LogEntry struct {
	Model
	Timestamp time.Time
	IP        string `gorm:"size:15"`
	Method    string `gorm:"size:10"`
	Path      string `gorm:"size:255"`
	Status    int    `gorm:"index"`
	LatencyMs int    `gorm:"column:latency_ms"`
}
