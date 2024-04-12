package initializer

import (
	"github.com/vishnusunil243/Job-Portal-Company-Service/concurrency"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/service"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/usecases"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.CompanyService {
	adapter := adapters.NewCompanyAdapter(db)
	usecases := usecases.NewCompanyUseCase(adapter)
	service := service.NewCompanyService(adapter, usecases, ":8087", ":8081", ":8083")
	c := concurrency.NewConcurrency(db)
	c.Concurrency()
	return service
}
