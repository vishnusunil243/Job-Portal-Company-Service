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
	GetAllJobForCompany(companyId string) ([]helperstruct.JobHelper, error)
	UpdateJob(string, helperstruct.JobHelper) error
	DeleteJob(ID string) error
	AddJobSkill(entities.JobSkill) error
	DeleteJobSkill(string) error
	GetAllJobSkills(jobId string) ([]entities.JobSkill, error)
	CreateProfile(entities.Profile) error
	AddLink(entities.Link) error
	DeleteLink(Id string) error
	GetAllLink(profileId string) ([]entities.Link, error)
	GetProfileIdFromCompanyId(companyId string) (string, error)
	GetCompanyById(companyId string) (entities.Company, error)
	AddAddress(entities.Address) error
	EditAddress(entities.Address) error
	GetAddress(profileId string) (entities.Address, error)
	EditName(entities.Company) error
	EditPhone(entities.Company) error
	UploadImage(profileId, image string) (string, error)
	GetProfilePic(string) (string, error)
	CompanyGetJobByDesignation(companyId, designation string) (entities.Job, error)
	CompanyGetJobSkill(jobId string, skillId int) (entities.JobSkill, error)
	JobSearch(designation, experience string, categoryId int) ([]helperstruct.JobHelper, error)
	GetHomeUsers(designation string) ([]helperstruct.JobHelper, error)
	NotifyMe(userId, companyId string) error
	GetNotifyMeByCompanyId(companyId string) ([]helperstruct.NotifyHelper, error)
	GetAllNotifyMe(userId string) ([]helperstruct.NotifyHelper, error)
	RemoveNotifyMe(userId, companyId string) error
	GetNotifyMe(companyId, userId string) (entities.NotifyMe, error)
	UpdateAverageRating(rating float64, companyId string) error
	GetAllCompanies() ([]entities.Company, error)
	BlockCompany(companyId string) error
	UnblockCompany(companyID string) error
	GetCompanyIdFromJobId(jobId string) (string, error)
	UpdateHired(jobId string) error
}
