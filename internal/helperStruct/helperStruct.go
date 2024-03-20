package helperstruct

import (
	"time"

	"github.com/google/uuid"
)

type JobHelper struct {
	JobID         uuid.UUID
	Designation   string
	Capacity      int
	Hired         int
	StatusID      int
	Status        string
	MinSalary     int64
	MaxSalary     int64
	MinExperience string
	Company       string
	PostedOn      time.Time
	ValidUntil    time.Time
}
type NotifyHelper struct {
	UserId    uuid.UUID
	CompanyId uuid.UUID
	UserEmail string
	Company   string
}
