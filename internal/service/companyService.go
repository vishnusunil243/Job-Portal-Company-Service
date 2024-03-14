package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/helper"
	helperstruct "github.com/vishnusunil243/Job-Portal-Company-Service/internal/helperStruct"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UserClient pb.UserServiceClient
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
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Phone:      req.Phone,
		CategoryId: int(req.CategoryId),
	}
	res, err := company.adapters.CompanySignup(reqEntity)
	return &pb.CompanySignupResponse{
		Id:         res.ID.String(),
		Name:       res.Name,
		Email:      res.Email,
		Phone:      res.Phone,
		CategoryId: int32(res.CategoryId),
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
func (company *CompanyService) AddJobs(ctx context.Context, req *pb.AddJobRequest) (*pb.JobResponse, error) {
	layout := "2006-01-02T15:04:05.999999Z"
	validUntil, err := time.Parse(layout, req.ValidUntil)
	if err != nil {
		return &pb.JobResponse{}, fmt.Errorf("please provide time in appropriate format")
	}
	companyID, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return &pb.JobResponse{}, err
	}
	jobreqEntity := entities.Job{
		Designation:   req.Designation,
		Capacity:      int(req.Vacancy),
		Hired:         0,
		ValidUntil:    validUntil,
		CompanyID:     companyID,
		MinExperience: req.MinExperience,
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
		PostedOn:    jobData.PostedOn.String(),
		ValidUntil:  jobData.ValidUntil.String(),
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
			PostedOn:    job.PostedOn.String(),
			ValidUntil:  job.ValidUntil.String(),
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
		PostedOn:    job.PostedOn.String(),
		ValidUntil:  job.ValidUntil.String(),
		Salaryrange: sRange,
	}, nil
}
func (company *CompanyService) GetAllJobsForCompany(req *pb.GetJobByCompanyId, srv pb.CompanyService_GetAllJobsForCompanyServer) error {
	jobs, err := company.adapters.GetAllJobForCompany(req.Id)
	if err != nil {
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
			PostedOn:    job.PostedOn.String(),
			ValidUntil:  job.ValidUntil.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) UpdateJobs(ctx context.Context, req *pb.UpdateJobRequest) (*emptypb.Empty, error) {
	var validUntil time.Time
	if req.ValidUntil != "" {
		layout := "2006-01-02T15:04:05.999999Z"
		time, err := time.Parse(layout, req.ValidUntil)
		validUntil = time
		if err != nil {
			return &emptypb.Empty{}, fmt.Errorf("please provide time in appropriate format")
		}
	}
	reqEntity := helperstruct.JobHelper{
		Designation: req.Designation,
		Capacity:    int(req.Capacity),
		Hired:       int(req.Hired),
		StatusID:    int(req.StatusId),
		MinSalary:   req.Salaryrange.MinSalary,
		MaxSalary:   req.Salaryrange.MaxSalary,
		ValidUntil:  validUntil,
	}
	err := company.adapters.UpdateJob(req.JobId, reqEntity)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (company *CompanyService) DeleteJob(ctx context.Context, req *pb.GetJobById) (*emptypb.Empty, error) {
	if req.Id == "" {
		return &emptypb.Empty{}, fmt.Errorf("invalid job id")
	}
	err := company.adapters.DeleteJob(req.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (company *CompanyService) CompanyAddJobSkill(ctx context.Context, req *pb.AddJobSkillRequest) (*emptypb.Empty, error) {
	skill, err := UserClient.GetSkillById(context.Background(), &pb.GetSkillByIdRequest{
		Id: req.SkillId,
	})
	if err != nil {
		return nil, fmt.Errorf("please enter a valid skill id")
	}
	if skill.Category == "" {
		return nil, fmt.Errorf("please enter a valid skill id")
	}
	jobId, err := uuid.Parse(req.JobId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.JobSkill{
		JobId:   jobId,
		SkillId: int(req.SkillId),
	}
	if err := company.adapters.AddJobSkill(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) DeleteJobSkill(ctx context.Context, req *pb.JobSkillId) (*emptypb.Empty, error) {
	if err := company.adapters.DeleteJobSkill(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) GetAllJobSkill(req *pb.GetJobById, srv pb.CompanyService_GetAllJobSkillServer) error {
	jobSkills, err := company.adapters.GetAllJobSkills(req.Id)
	if err != nil {
		return err
	}
	for _, jobSkill := range jobSkills {
		skill, err := UserClient.GetSkillById(context.Background(), &pb.GetSkillByIdRequest{
			Id: int32(jobSkill.SkillId),
		})
		if err != nil {
			return err
		}
		res := &pb.JobSkillResponse{
			Id:      jobSkill.ID.String(),
			SkillId: int32(jobSkill.SkillId),
			Skill:   skill.Skill,
			JobId:   jobSkill.JobId.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) CompanyAddLink(ctx context.Context, req *pb.CompanyLinkRequest) (*emptypb.Empty, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.CompanyId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Link{
		Title:     req.Title,
		URL:       req.Url,
		ProfileId: profileId,
	}
	if err := company.adapters.AddLink(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) CompanyDeleteLink(ctx context.Context, req *pb.CompanyDeleteLinkRequest) (*emptypb.Empty, error) {
	if err := company.adapters.DeleteLink(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) CompanyGetAllLink(req *pb.GetJobByCompanyId, srv pb.CompanyService_CompanyGetAllLinkServer) error {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.Id)
	if err != nil {
		return err
	}

	links, err := company.adapters.GetAllLink(profile)
	if err != nil {
		return err
	}
	for _, link := range links {
		res := &pb.CompanyLinkResponse{
			Id:    link.ID.String(),
			Title: link.Title,
			Url:   link.URL,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) CompanyCreateProfile(ctx context.Context, req *pb.GetJobByCompanyId) (*emptypb.Empty, error) {
	companyId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Profile{
		CompanyId: companyId,
	}
	if err := company.adapters.CreateProfile(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) GetCompanyById(ctx context.Context, req *pb.GetJobByCompanyId) (*pb.CompanySignupResponse, error) {
	cmpny, err := company.adapters.GetCompanyById(req.Id)
	if err != nil {
		return &pb.CompanySignupResponse{}, err
	}
	res := &pb.CompanySignupResponse{
		Id:         cmpny.ID.String(),
		Email:      cmpny.Email,
		Name:       cmpny.Name,
		Phone:      cmpny.Phone,
		CategoryId: int32(cmpny.CategoryId),
	}
	return res, nil
}
func (company *CompanyService) CompanyAddAddress(ctx context.Context, req *pb.CompanyAddAddressRequest) (*emptypb.Empty, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.CompanyId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	address, err := company.adapters.GetAddress(profile)
	if err != nil {
		return nil, err
	}
	if address.Country != "" {
		return nil, fmt.Errorf("address already exist please edit  the existing one")
	}
	reqEntity := entities.Address{
		ProfileId: profileId,
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
	}
	if err := company.adapters.AddAddress(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) CompanyEditAddress(ctx context.Context, req *pb.CompanyAddAddressRequest) (*emptypb.Empty, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.CompanyId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	addr, err := company.adapters.GetAddress(profile)
	if err != nil {
		return nil, err
	}
	if addr.Country == "" {
		return nil, fmt.Errorf("please add an address first")
	}
	reqEntity := entities.Address{
		ProfileId: profileId,
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
	}
	if err := company.adapters.EditAddress(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) CompanyGetAddress(ctx context.Context, req *pb.GetJobByCompanyId) (*pb.CompanyAddressResponse, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.Id)
	if err != nil {
		return nil, err
	}
	address, err := company.adapters.GetAddress(profile)
	if err != nil {
		return nil, err
	}
	res := &pb.CompanyAddressResponse{
		Country:  address.Country,
		State:    address.State,
		District: address.District,
		City:     address.City,
		Id:       address.ID.String(),
	}
	return res, nil
}
func (company *CompanyService) CompanyEditName(ctx context.Context, req *pb.CompanyEditNameRequest) (*emptypb.Empty, error) {
	companyId, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Company{
		Name: req.Name,
		ID:   companyId,
	}
	if err := company.adapters.EditName(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) CompanyEditPhone(ctx context.Context, req *pb.CompanyEditPhoneRequest) (*emptypb.Empty, error) {
	companyId, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Company{
		Phone: req.Phone,
		ID:    companyId,
	}
	if err := company.adapters.EditPhone(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
