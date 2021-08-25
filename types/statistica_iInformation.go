package types

import "time"

type StatisticalInformation struct {
	Key      string
	Value    string
	TimeUsed time.Duration
}
