package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (company *CompanyService) AddJob(ctx context.Context, req *pb.AddJobRequest) (*pb.JobResponse, error) {
	validUntil, err := ptypes.Timestamp(req.ValidUntil)
	if err != nil {
		return nil, err
	}
	jobreqEntity := entities.Job{
		Designation: req.Designation,
		Capacity:    int(req.Vacancy),
		Hired:       0,
		ValidUntil:  validUntil,
	}
	salaryRangeEntity := entities.SalaryRange{
		MinSalary: req.Salaryrange.MinSalary,
		MaxSalary: req.Salaryrange.MaxSalary,
	}
	jobData, sRange, err := company.adapters.AddJob(jobreqEntity, salaryRangeEntity)
	if err != nil {
		return &pb.JobResponse{}, err
	}
	resSalaryRange := pb.SalaryRange{
		MinSalary: sRange.MinSalary,
		MaxSalary: sRange.MaxSalary,
	}
	return &pb.JobResponse{
		Designation: jobData.Designation,
		PostedOn:    timestamppb.New(jobData.PostedOn),
		ValidUntil:  timestamppb.New(jobData.ValidUntil),
		Vacancy:     int32(jobData.Capacity),
		Hired:       int32(jobData.Hired),
		Salaryrange: &resSalaryRange,
	}, nil
}
func (company *CompanyService) GetAllJobs(e *emptypb.Empty, srv pb.CompanyService_GetAllJobsServer) error {
	jobs, err := company.adapters.GetAllJobs()
	if err != nil {
		fmt.Println("error fetching jobs ", err.Error())
		return err
	}
	for _, job := range jobs {
		resSalaryRange := pb.SalaryRange{
			MinSalary: job.MinSalary,
			MaxSalary: job.MaxSalary,
		}
		res := &pb.JobResponse{
			Designation: job.Designation,
			Salaryrange: &resSalaryRange,
			Vacancy:     int32(job.Capacity),
			Hired:       int32(job.Hired),
			PostedOn:    timestamppb.New(job.PostedOn),
			ValidUntil:  timestamppb.New(job.ValidUntil),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) GetJob(ctx context.Context, req *pb.GetJobById) (*pb.JobResponse, error) {
	job, err := company.adapters.GetJob(req.Id)
	if err != nil {
		return &pb.JobResponse{}, err
	}
	sRange := &pb.SalaryRange{
		MinSalary: job.MinSalary,
		MaxSalary: job.MaxSalary,
	}
	return &pb.JobResponse{
		Designation: job.Designation,
		Vacancy:     int32(job.Capacity),
		Hired:       int32(job.Hired),
		PostedOn:    timestamppb.New(job.PostedOn),
		ValidUntil:  timestamppb.New(job.ValidUntil),
		Salaryrange: sRange,
	}, nil
}
