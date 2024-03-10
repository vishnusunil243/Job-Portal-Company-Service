package adapters

import (
	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	helperstruct "github.com/vishnusunil243/Job-Portal-Company-Service/internal/helperStruct"
	"gorm.io/gorm"
)

type CompanyAdapter struct {
	DB *gorm.DB
}

func NewCompanyAdapter(db *gorm.DB) *CompanyAdapter {
	return &CompanyAdapter{
		DB: db,
	}
}
func (company *CompanyAdapter) CompanySignup(req entities.Company) (entities.Company, error) {
	id := uuid.New()
	var res entities.Company
	insertQuery := `INSERT INTO company (id,name,email,phone,password) VALUES ($1,$2,$3,$4,$5) RETURNING *`
	if err := company.DB.Raw(insertQuery, id, req.Name, req.Email, req.Phone, req.Password).Scan(&res).Error; err != nil {
		return entities.Company{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetCompanyByEmail(email string) (entities.Company, error) {
	selectQuery := `SELECT * FROM company WHERE email=?`
	var res entities.Company
	if err := company.DB.Raw(selectQuery, email).Scan(&res).Error; err != nil {
		return entities.Company{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) AddJob(jobreq entities.Job, salaryRange entities.SalaryRange) (entities.Job, entities.SalaryRange, error) {
	tx := company.DB.Begin()
	jobId := uuid.New()
	var jobData entities.Job
	insertJobQuery := `INSERT INTO jobs (id,deisgnation,capacity,hired,status_id,posted_on,valid_until) VALUES ($1,$2,$3,$4,$5,NOW(),$6) RETURNING *`
	if err := tx.Raw(insertJobQuery, jobId, jobreq.Designation, jobreq.Capacity, jobreq.Hired, 1, jobreq.ValidUntil).Scan(&jobData).Error; err != nil {
		tx.Rollback()
		return entities.Job{}, entities.SalaryRange{}, err
	}
	var sRange entities.SalaryRange
	if salaryRange.MaxSalary != 0 && salaryRange.MinSalary != 0 {
		salaryId := uuid.New()
		insertSalaryRangeQuery := `INSERT INTO salary_ranges (id,max_salary,min_salary,job_id) VALUES ($1,$2,$3,$4) RETURNING *`
		if err := tx.Raw(insertSalaryRangeQuery, salaryId, salaryRange.MaxSalary, salaryRange.MinSalary, jobId).Scan(&sRange).Error; err != nil {
			tx.Rollback()
			return entities.Job{}, entities.SalaryRange{}, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return entities.Job{}, entities.SalaryRange{}, err
	}
	return jobData, sRange, nil
}
func (company *CompanyAdapter) GetAllJobs() ([]helperstruct.JobHelper, error) {
	selectQuery := `SELECT j.id AS job_id,designation,capacity,hired,status,max_salary,min_salary FROM jobs j LEFT JOIN salary_ranges s`
	var res []helperstruct.JobHelper
	if err := company.DB.Raw(selectQuery).Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetJob(ID string) (helperstruct.JobHelper, error) {
	selectQuery := `SELECT j.id AS job_id,designation,capacity,hired,status,max_salary,min_salary FROM jobs j LEFT JOIN salary_ranges s WHERE id=?`
	var res helperstruct.JobHelper
	if err := company.DB.Raw(selectQuery, ID).Scan(&res).Error; err != nil {
		return helperstruct.JobHelper{}, err
	}
	return res, nil

}
