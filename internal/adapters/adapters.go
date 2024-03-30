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
	insertQuery := `INSERT INTO companies (id,name,email,phone,password,category_id,avg_rating,created_at) VALUES ($1,$2,$3,$4,$5,$6,0,NOW()) RETURNING *`
	if err := company.DB.Raw(insertQuery, id, req.Name, req.Email, req.Phone, req.Password, req.CategoryId).Scan(&res).Error; err != nil {
		return entities.Company{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetCompanyByEmail(email string) (entities.Company, error) {
	selectQuery := `SELECT * FROM companies WHERE email=?`
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
	insertJobQuery := `INSERT INTO jobs (id,designation,capacity,hired,status_id,posted_on,valid_until,company_id,min_experience) VALUES ($1,$2,$3,$4,$5,NOW(),$6,$7,$8) RETURNING *`
	if err := tx.Raw(insertJobQuery, jobId, jobreq.Designation, jobreq.Capacity, jobreq.Hired, 1, jobreq.ValidUntil, jobreq.CompanyID, jobreq.MinExperience).Scan(&jobData).Error; err != nil {
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
	selectQuery := `SELECT j.id AS job_id,designation,capacity,hired,status,max_salary,min_salary,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id`
	var res []helperstruct.JobHelper
	if err := company.DB.Raw(selectQuery).Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetJob(ID string) (helperstruct.JobHelper, error) {
	selectQuery := `SELECT j.id AS job_id,designation,capacity,hired,status,max_salary,min_salary,posted_on,valid_until,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id WHERE j.id=?`
	var res helperstruct.JobHelper
	if err := company.DB.Raw(selectQuery, ID).Scan(&res).Error; err != nil {
		return helperstruct.JobHelper{}, err
	}
	return res, nil

}
func (company *CompanyAdapter) GetAllJobForCompany(companyId string) ([]helperstruct.JobHelper, error) {
	var res []helperstruct.JobHelper
	selectQuery := `SELECT j.id AS job_id,max_salary,min_salary,designation,valid_until,posted_on,company_id,capacity,hired,status,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id WHERE company_id=?`
	if err := company.DB.Raw(selectQuery, companyId).Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) UpdateJob(ID string, req helperstruct.JobHelper) error {
	updateJobs := `UPDATE jobs SET designation=$1,capacity=$2,hired=$3,status_id=$4,valid_until=$5,min_experience=$6 WHERE id=$7`
	tx := company.DB.Begin()
	if err := tx.Exec(updateJobs, req.Designation, req.Capacity, req.Hired, req.StatusID, req.ValidUntil, req.MinExperience, ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	updateSalaryRange := `UPDATE salary_ranges SET min_salary=$1,max_salary=$2 WHERE job_id=$3`
	if err := tx.Exec(updateSalaryRange, req.MinSalary, req.MaxSalary, ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) DeleteJob(ID string) error {
	deleteQuery := `DELETE FROM jobs where id=?`
	if err := company.DB.Exec(deleteQuery, ID).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) AddJobSkill(req entities.JobSkill) error {
	addJobSkillQuery := `INSERT INTO job_skills(id,skill_id,job_id) VALUES ($1,$2,$3)`
	id := uuid.New()
	if err := company.DB.Exec(addJobSkillQuery, id, req.SkillId, req.JobId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) DeleteJobSkill(id string) error {
	updateJobSkillQuery := `DELETE FROM job_skills WHERE id=?`
	if err := company.DB.Exec(updateJobSkillQuery, id).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetAllJobSkills(jobId string) ([]entities.JobSkill, error) {
	getAllJobSkillQuery := `SELECT * FROM job_skills WHERE job_id=?`
	var res []entities.JobSkill
	if err := company.DB.Raw(getAllJobSkillQuery, jobId).Scan(&res).Error; err != nil {
		return []entities.JobSkill{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) CreateProfile(req entities.Profile) error {
	id := uuid.New()
	insertProfileQuery := `INSERT INTO profiles (id,company_id) VALUES ($1,$2)`
	if err := company.DB.Exec(insertProfileQuery, id, req.CompanyId).Error; err != nil {
		return err
	}
	return nil

}
func (company *CompanyAdapter) AddLink(req entities.Link) error {
	id := uuid.New()
	insertLinkQuery := `INSERT INTO links (id,url,title,profile_id) VALUES ($1,$2,$3,$4)`
	if err := company.DB.Exec(insertLinkQuery, id, req.URL, req.Title, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) DeleteLink(id string) error {
	deleteLinkQuery := `DELETE FROM links WHERE id=?`
	if err := company.DB.Exec(deleteLinkQuery, id).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetAllLink(profileId string) ([]entities.Link, error) {
	var res []entities.Link
	selectLinks := `SELECT * FROM links WHERE profile_id=?`
	if err := company.DB.Raw(selectLinks, profileId).Scan(&res).Error; err != nil {
		return []entities.Link{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetProfileIdFromCompanyId(companyId string) (string, error) {
	var res string
	selectProfileQuery := `SELECT id FROM profiles WHERE company_id=?`
	if err := company.DB.Raw(selectProfileQuery, companyId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}
func (company *CompanyAdapter) GetCompanyById(companyId string) (entities.Company, error) {
	selectCompanyQuery := `SELECT * FROM companies WHERE id=?`
	var res entities.Company
	if err := company.DB.Raw(selectCompanyQuery, companyId).Scan(&res).Error; err != nil {
		return entities.Company{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) EditName(req entities.Company) error {
	updateQuery := `UPDATE companies SET name=$1 WHERE id=$2`
	if err := company.DB.Exec(updateQuery, req.Name, req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) EditPhone(req entities.Company) error {
	updateQuery := `UPDATE companies SET phone=$1 WHERE id=$2`
	if err := company.DB.Exec(updateQuery, req.Phone, req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) AddAddress(req entities.Address) error {
	id := uuid.New()

	insertQuery := `INSERT INTO addresses (id,country,state,district,city,profile_id) VALUES ($1,$2,$3,$4,$5,$6)`
	if err := company.DB.Exec(insertQuery, id, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) EditAddress(req entities.Address) error {
	updateQuery := `UPDATE addresses SET country=$1,state=$2,district=$3,city=$4 WHERE profile_id=$5`
	if err := company.DB.Exec(updateQuery, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetAddress(profileId string) (entities.Address, error) {
	selectQuery := `SELECT * FROM addresses WHERE profile_id=?`
	var res entities.Address
	if err := company.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return entities.Address{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetCompanyIdFromJobId(jobId string) (string, error) {
	var companyId string
	selectQuery := `SELECT company_id FROM jobs WHERE id=?`
	if err := company.DB.Raw(selectQuery, jobId).Scan(&companyId).Error; err != nil {
		return "", err
	}
	return companyId, nil
}
func (company *CompanyAdapter) UploadImage(image, profileId string) (string, error) {
	var res string
	insertImageQuery := `UPDATE profiles SET image=$1 WHERE id=$2 RETURNING image`
	if err := company.DB.Raw(insertImageQuery, image, profileId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}
func (company *CompanyAdapter) GetProfilePic(profileId string) (string, error) {
	selectQuery := `SELECT image FROM profiles WHERE id=?`
	var image string
	if err := company.DB.Raw(selectQuery, profileId).Scan(&image).Error; err != nil {
		return "", nil
	}
	return image, nil
}
func (company *CompanyAdapter) CompanyGetJobByDesignation(companyId, designation string) (entities.Job, error) {
	selectQuery := `SELECT * FROM jobs WHERE company_id=$1 AND designation=$2`
	var res entities.Job
	if err := company.DB.Raw(selectQuery, companyId, designation).Scan(&res).Error; err != nil {
		return entities.Job{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) CompanyGetJobSkill(jobId string, skillId int) (entities.JobSkill, error) {

	var res entities.JobSkill
	selectSkillQuery := `SELECT * FROM job_skills WHERE job_id=$1 AND skill_id=$2`
	if err := company.DB.Raw(selectSkillQuery, jobId, skillId).Scan(&res).Error; err != nil {
		return entities.JobSkill{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) JobSearch(designation, experience string) ([]helperstruct.JobHelper, error) {
	selectJobQuery := `SELECT j.id AS job_id,max_salary,min_salary,designation,valid_until,posted_on,company_id,capacity,hired,status,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id WHERE designation ILIKE $1 ORDER BY posted_on DESC`
	var res []helperstruct.JobHelper
	if err := company.DB.Raw(selectJobQuery, "%"+designation+"%").Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetHomeUsers(designation string) ([]helperstruct.JobHelper, error) {
	selectJobQuery := `WITH senior_jobs AS (
	SELECT j.id AS job_id,max_salary,min_salary,designation,valid_until,posted_on,company_id,capacity,hired,status,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id WHERE designation ILIKE $1 ORDER BY posted_on DESC
	)
	SELECT * FROM senior_jobs
	UNION ALL
	SELECT j.id AS job_id,max_salary,min_salary,designation,valid_until,posted_on,company_id,capacity,hired,status,c.name AS company,min_experience FROM jobs j LEFT JOIN salary_ranges s ON s.job_id=j.id LEFT JOIN statuses ON j.status_id=statuses.id LEFT JOIN companies c ON c.id=j.company_id WHERE designation NOT ILIKE $1`
	var res []helperstruct.JobHelper
	if err := company.DB.Raw(selectJobQuery, "%"+designation+"%").Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) NotifyMe(userId, companyId string) error {
	id := uuid.New()
	insertIntoNotifyMeQuery := `INSERT INTO notify_mes (id,company_id,user_id) VALUES ($1,$2,$3)`
	if err := company.DB.Exec(insertIntoNotifyMeQuery, id, companyId, userId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetNotifyMeByCompanyId(companyId string) ([]helperstruct.NotifyHelper, error) {
	selectQuery := `SELECT n.user_id,n.company_id,c.name AS company FROM notify_mes n JOIN companies c ON n.company_id=c.id  WHERE n.company_id=?`
	var res []helperstruct.NotifyHelper
	if err := company.DB.Raw(selectQuery, companyId).Scan(&res).Error; err != nil {
		return []helperstruct.NotifyHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) GetAllNotifyMe(userId string) ([]helperstruct.NotifyHelper, error) {
	selectQuery := `SELECT n.user_id,n.company_id,c.name AS company FROM notify_mes n JOIN companies c ON n.company_id=c.id WHERE n.user_id=?`
	var res []helperstruct.NotifyHelper
	if err := company.DB.Raw(selectQuery, userId).Scan(&res).Error; err != nil {
		return []helperstruct.NotifyHelper{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) RemoveNotifyMe(userId, companyId string) error {
	deleteQuery := `DELETE FROM notify_mes WHERE user_id=$1 AND company_id=$2`
	if err := company.DB.Exec(deleteQuery, userId, companyId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetNotifyMe(companyId, userId string) (entities.NotifyMe, error) {
	var res entities.NotifyMe
	selectQuery := `SELECT * FROM notify_mes WHERE company_id=$1 AND user_id=$2`
	if err := company.DB.Raw(selectQuery, companyId, userId).Scan(&res).Error; err != nil {
		return entities.NotifyMe{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) UpdateAverageRating(rating float64, companyId string) error {
	updateRatingQuery := `UPDATE companies SET avg_rating=$1 WHERE id=$2`
	if err := company.DB.Exec(updateRatingQuery, rating, companyId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) GetAllCompanies() ([]entities.Company, error) {
	selectQuery := `SELECT * FROM companies ORDER BY avg_rating DESC`
	var res []entities.Company
	if err := company.DB.Raw(selectQuery).Scan(&res).Error; err != nil {
		return []entities.Company{}, err
	}
	return res, nil
}
func (company *CompanyAdapter) BlockCompany(companyId string) error {
	updateQuery := `UPDATE companies SET is_blocked=true WHERE id=?`
	if err := company.DB.Exec(updateQuery, companyId).Error; err != nil {
		return err
	}
	return nil
}
func (company *CompanyAdapter) UnblockCompany(companyID string) error {
	updateQuery := `UPDATE companies SET is_blocked=false WHERE id=?`
	if err := company.DB.Exec(updateQuery, companyID).Error; err != nil {
		return err
	}
	return nil
}
