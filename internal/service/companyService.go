package service

import (
	"context"
	"fmt"

	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type CompanyService struct {
	adapters adapters.AdapterInterface
	pb.UnimplementedCompanyServiceServer
}

func NewCompanyService(adapters adapters.AdapterInterface) *CompanyService {
	return &CompanyService{
		adapters: adapters,
	}
}
func (company CompanyService) CompanySignup(ctx context.Context, req *pb.CompanySignupRequest) (*pb.CompanySignupResponse, error) {
	if req.Name == "" {
		return &pb.CompanySignupResponse{}, fmt.Errorf("please enter a valid name")
	}
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return &pb.CompanySignupResponse{}, err
	}
	reqEntity := entities.Company{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
	}
	res, err := company.adapters.CompanySignup(reqEntity)
	return &pb.CompanySignupResponse{
		Id:    res.ID.String(),
		Name:  res.Name,
		Email: res.Email,
		Phone: res.Phone,
	}, err
}
func (company *CompanyService) CompanyLogin(ctx context.Context, req *pb.CompanyLoginRequest) (*pb.CompanySignupResponse, error) {
	if req.Email == "" {
		return &pb.CompanySignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	companyData, err := company.adapters.GetCompanyByEmail(req.Email)
	if err != nil {
		return &pb.CompanySignupResponse{}, err
	}
	if companyData.Email == "" {
		return &pb.CompanySignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(companyData.Password, req.Password) {
		return &pb.CompanySignupResponse{}, fmt.Errorf("invalid credentials")
	}
	return &pb.CompanySignupResponse{
		Id:    companyData.ID.String(),
		Email: companyData.Email,
		Name:  companyData.Name,
		Phone: companyData.Phone,
	}, nil
}
