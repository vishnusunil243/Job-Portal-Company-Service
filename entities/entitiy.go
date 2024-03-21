package entities

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID         uuid.UUID
	Name       string
	Email      string
	Phone      string
	CategoryId int
	Password   string
	AvgRating  float64
	IsBlocked  bool
}

type Address struct {
	ID        uuid.UUID
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	Country   string
	State     string
	District  string
	City      string
}
type Link struct {
	ID        uuid.UUID
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	Title     string
	URL       string
}
type Job struct {
	ID            uuid.UUID
	CompanyID     uuid.UUID
	Designation   string
	Capacity      int
	Hired         int
	StatusId      int
	Status        Status `gorm:"foreignKey:StatusId"`
	MinExperience string
	PostedOn      time.Time
	ValidUntil    time.Time
}
type JobSkill struct {
	ID      uuid.UUID
	JobId   uuid.UUID
	Job     Job `gorm:"foreignKey:JobId"`
	SkillId int
}
type SalaryRange struct {
	ID        uuid.UUID
	JobID     uuid.UUID
	Job       Job `gorm:"foreignKey:JobID"`
	MinSalary int64
	MaxSalary int64
}
type Status struct {
	ID     int
	Status string
}
type Profile struct {
	ID        uuid.UUID
	CompanyId uuid.UUID
	Company   Company `gorm:"foreignKey:CompanyId"`
	Image     string
}
type NotifyMe struct {
	ID        uuid.UUID
	CompanyId uuid.UUID
	Company   Company `gorm:"foreignKey:CompanyId"`
	UserId    uuid.UUID
}
