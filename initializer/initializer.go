package initializer

import (
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.CompanyService {
	adapter := adapters.NewCompanyAdapter(db)
	service := service.NewCompanyService(adapter)
	return service
}
