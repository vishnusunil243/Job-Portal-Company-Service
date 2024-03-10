package helperstruct

import (
	"time"

	"github.com/google/uuid"
)

type JobHelper struct {
	JobID       uuid.UUID
	Designation string
	Capacity    int
	Hired       int
	StatusID    int
	Status      string
	MinSalary   int64
	MaxSalary   int64
	PostedOn    time.Time
	ValidUntil  time.Time
}
