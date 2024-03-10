package adapters

import (
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	helperstruct "github.com/vishnusunil243/Job-Portal-Company-Service/internal/helperStruct"
)

type AdapterInterface interface {
	CompanySignup(entities.Company) (entities.Company, error)
	GetCompanyByEmail(string) (entities.Company, error)
	AddJob(entities.Job, entities.SalaryRange) (entities.Job, entities.SalaryRange, error)
	GetAllJobs() ([]helperstruct.JobHelper, error)
	GetJob(ID string) (helperstruct.JobHelper, error)
}
