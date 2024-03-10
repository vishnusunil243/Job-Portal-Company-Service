package adapters

import (
	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
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
