package models

import "time"

type SignalJobContext struct {
	Execute  bool
	FromTime time.Time
	DateKey  string
}
