package adapters

import "github.com/vishnusunil243/Job-Portal-Company-Service/entities"

type AdapterInterface interface {
	CompanySignup(entities.Company) (entities.Company, error)
	GetCompanyByEmail(string) (entities.Company, error)
}
