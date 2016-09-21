package models

import "time"

type Plan struct {
	Plan_id        int `json:"plan_id, omitempty"`
	Plan_number    string
	Plan_type      string
	Specification1 string
	Specification2 string
	Price          float32
	Cycle          string
	Create_time    time.Time
	Status         string
}
