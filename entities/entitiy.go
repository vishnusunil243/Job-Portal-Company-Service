package entities

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Phone    string
	Password string
}
type Address struct {
	ID       uuid.UUID
	Country  string
	State    string
	District string
	City     string
}
type Link struct {
	ID    uuid.UUID
	Title string
	URL   string
}
type Job struct {
	ID          uuid.UUID
	Designation string
	Capacity    int
	Hired       int
	StatusId    int
	Status      Status `gorm:"foreignKey:StatusId"`
	PostedOn    time.Time
	ValidUntil  time.Time
}
type SalaryRange struct {
	ID        uuid.UUID
	JobID     uuid.UUID
	Job       Job `gorm:"foreignKey:jobID"`
	MinSalary int64
	MaxSalary int64
}
type Status struct {
	ID     int
	Status string
}
