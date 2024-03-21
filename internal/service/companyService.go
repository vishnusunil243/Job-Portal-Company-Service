package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/helper"
	helperstruct "github.com/vishnusunil243/Job-Portal-Company-Service/internal/helperStruct"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/usecases"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UserClient         pb.UserServiceClient
	SearchClient       pb.SearchServiceClient
	NotificationClient pb.EmailServiceClient
)

type CompanyService struct {
	adapters adapters.AdapterInterface
	usecases usecases.Usecase
	pb.UnimplementedCompanyServiceServer
}

func NewCompanyService(adapters adapters.AdapterInterface, usecases usecases.Usecase) *CompanyService {
	return &CompanyService{
		adapters: adapters,
		usecases: usecases,
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
	check, err := company.adapters.GetCompanyByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if check.Name != "" {
		return nil, fmt.Errorf("a company is already registered with the given email id")
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
	if companyData.IsBlocked {
		return &pb.CompanySignupResponse{}, fmt.Errorf("unfortunately you have been blocked by the admin")
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
	job, err := company.adapters.CompanyGetJobByDesignation(req.CompanyId, req.Designation)
	if err != nil {
		return nil, err
	}
	notifyme, err := company.adapters.GetNotifyMeByCompanyId(req.CompanyId)
	if err != nil {
		return nil, err
	}
	go func(notifyme []helperstruct.NotifyHelper) {
		for _, ntify := range notifyme {
			rq := &pb.AddNotificationRequest{
				UserId:  ntify.UserId.String(),
				Message: fmt.Sprintf(`{"message":"%s has posted a new job opening , check it out "}`, ntify.Company),
			}
			_, err := NotificationClient.AddNotification(context.Background(), rq)
			if err != nil {
				log.Printf("notification not sent err:%v", err)
			} else {
				log.Printf("notification sent to %s", ntify.UserId)
			}
		}
	}(notifyme)
	if job.Designation != "" {
		return nil, fmt.Errorf("you have already added a job for the given designation please add a new designation or update the previous job post")
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
	job, sRange, err := company.adapters.AddJob(jobreqEntity, salaryRangeEntity)
	if err != nil {
		return &pb.JobResponse{}, err
	}
	resSalaryRange := pb.SalaryRange{
		MinSalary: sRange.MinSalary,
		MaxSalary: sRange.MaxSalary,
	}
	return &pb.JobResponse{
		Designation:   job.Designation,
		Salaryrange:   &resSalaryRange,
		Vacancy:       int32(job.Capacity),
		Hired:         int32(job.Hired),
		PostedOn:      job.PostedOn.String(),
		ValidUntil:    job.ValidUntil.String(),
		Minexperience: job.MinExperience,
		Capacity:      int32(job.Capacity),
		Id:            job.ID.String(),
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
			Designation:   job.Designation,
			Salaryrange:   &resSalaryRange,
			Vacancy:       int32(job.Capacity) - int32(job.Hired),
			Hired:         int32(job.Hired),
			PostedOn:      job.PostedOn.String(),
			ValidUntil:    job.ValidUntil.String(),
			Company:       job.Company,
			Minexperience: job.MinExperience,
			Status:        job.Status,
			Capacity:      int32(job.Capacity),
			Id:            job.JobID.String(),
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
		Designation:   job.Designation,
		Salaryrange:   sRange,
		Vacancy:       int32(job.Capacity) - int32(job.Hired),
		Hired:         int32(job.Hired),
		PostedOn:      job.PostedOn.String(),
		ValidUntil:    job.ValidUntil.String(),
		Company:       job.Company,
		Minexperience: job.MinExperience,
		Status:        job.Status,
		Capacity:      int32(job.Capacity),
		Id:            job.JobID.String(),
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
			Designation:   job.Designation,
			Salaryrange:   &resSalaryRange,
			Vacancy:       int32(job.Capacity) - int32(job.Hired),
			Hired:         int32(job.Hired),
			PostedOn:      job.PostedOn.String(),
			ValidUntil:    job.ValidUntil.String(),
			Company:       job.Company,
			Minexperience: job.MinExperience,
			Status:        job.Status,
			Capacity:      int32(job.Capacity),
			Id:            job.JobID.String(),
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
		Designation:   req.Designation,
		Capacity:      int(req.Capacity),
		Hired:         int(req.Hired),
		StatusID:      int(req.StatusId),
		MinSalary:     req.Salaryrange.MinSalary,
		MaxSalary:     req.Salaryrange.MaxSalary,
		ValidUntil:    validUntil,
		MinExperience: req.MinExperience,
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
	jobSkill, err := company.adapters.CompanyGetJobSkill(req.JobId, int(req.SkillId))
	if err != nil {
		return nil, err
	}
	if jobSkill.SkillId != 0 {
		return nil, fmt.Errorf("this skill is aleady added please add a new skill")
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
	addressId := ""
	if address.ID != uuid.Nil {
		addressId = address.ID.String()
	}
	res := &pb.CompanyAddressResponse{
		Country:  address.Country,
		State:    address.State,
		District: address.District,
		City:     address.City,
		Id:       addressId,
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
func (company *CompanyService) CompanyUploadProfileImage(ctx context.Context, req *pb.CompanyImageRequest) (*pb.CompanyImageResponse, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.CompanyId)
	if err != nil {
		return nil, err
	}
	url, err := company.usecases.UploadImage(req, profile)
	if err != nil {
		return nil, err
	}
	return &pb.CompanyImageResponse{
		Url: url,
	}, nil

}
func (company *CompanyService) GetProfilePic(ctx context.Context, req *pb.GetJobByCompanyId) (*pb.CompanyImageResponse, error) {
	profile, err := company.adapters.GetProfileIdFromCompanyId(req.Id)
	if err != nil {
		return nil, err
	}
	image, err := company.adapters.GetProfilePic(profile)
	if err != nil {
		return nil, err
	}
	return &pb.CompanyImageResponse{
		Url: image,
	}, nil
}
func (company *CompanyService) JobSearch(req *pb.JobSearchRequest, srv pb.CompanyService_JobSearchServer) error {
	jobs, err := company.adapters.JobSearch(req.Designation, "")
	if err != nil {
		return err
	}
	for _, job := range jobs {
		salary := pb.SalaryRange{
			MinSalary: job.MinSalary,
			MaxSalary: job.MaxSalary,
		}
		res := &pb.JobResponse{
			Id:            job.JobID.String(),
			Company:       job.Company,
			Designation:   job.Designation,
			Salaryrange:   &salary,
			Vacancy:       int32(job.Capacity) - int32(job.Hired),
			Hired:         int32(job.Hired),
			Capacity:      int32(job.Capacity),
			PostedOn:      job.PostedOn.String(),
			ValidUntil:    job.ValidUntil.String(),
			Status:        job.Status,
			Minexperience: job.MinExperience,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}

	_, err = SearchClient.AddSearchHistory(context.Background(), &pb.SearchRequest{
		UserId:  req.UserId,
		Keyword: req.Designation,
	})
	if err != nil {
		return err
	}

	return nil
}
func (company *CompanyService) GetHome(req *pb.GetHomeRequest, srv pb.CompanyService_GetHomeServer) error {
	searchHistory, err := SearchClient.GetSearchHistory(context.Background(), &pb.UserId{
		UserId: req.UserId,
	})
	if err != nil {
		return err
	}
	jobs, err := company.adapters.GetHomeUsers(searchHistory.Designation)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		sRange := &pb.SalaryRange{
			MinSalary: job.MinSalary,
			MaxSalary: job.MaxSalary,
		}
		res := &pb.JobResponse{
			Id:            job.JobID.String(),
			Company:       job.Company,
			Designation:   job.Designation,
			Salaryrange:   sRange,
			Vacancy:       int32(job.Capacity) - int32(job.Hired),
			Hired:         int32(job.Hired),
			Capacity:      int32(job.Capacity),
			PostedOn:      job.PostedOn.String(),
			ValidUntil:    job.ValidUntil.String(),
			Status:        job.Status,
			Minexperience: job.MinExperience,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) NotifyMe(ctx context.Context, req *pb.NotifyMeRequest) (*emptypb.Empty, error) {
	check, err := company.adapters.GetNotifyMe(req.CompanyId, req.UserId)
	if err != nil {
		return nil, err
	}
	if check.CompanyId != uuid.Nil {
		return nil, fmt.Errorf("notifications are already enabled")
	}
	if err := company.adapters.NotifyMe(req.UserId, req.CompanyId); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) GetAllNotifyMe(req *pb.GetHomeRequest, srv pb.CompanyService_GetAllNotifyMeServer) error {
	companies, err := company.adapters.GetAllNotifyMe(req.UserId)
	if err != nil {
		return err
	}
	for _, company := range companies {
		res := &pb.NotifyMeResponse{
			Company: company.Company,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) CancelNotify(ctx context.Context, req *pb.NotifyMeRequest) (*emptypb.Empty, error) {
	if err := company.adapters.RemoveNotifyMe(req.UserId, req.CompanyId); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) UpdateAverageRatingOfCompany(ctx context.Context, req *pb.UpdateRatingRequest) (*emptypb.Empty, error) {
	if err := company.adapters.UpdateAverageRating(float64(req.AvgRating), req.CompanyId); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) GetAllCompany(e *emptypb.Empty, srv pb.CompanyService_GetAllCompanyServer) error {
	companies, err := company.adapters.GetAllCompanies()
	if err != nil {
		return err
	}

	for _, company := range companies {
		category, err := UserClient.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
			Id: int32(company.CategoryId),
		})
		if err != nil {
			return err
		}
		res := &pb.CompanyResponse{
			Id:        company.ID.String(),
			Name:      company.Name,
			Email:     company.Email,
			AvgRating: float32(company.AvgRating),
			Phone:     company.Phone,
			Category:  category.Category,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (company *CompanyService) GetCompany(ctx context.Context, req *pb.GetJobByCompanyId) (*pb.GetCompanyResponse, error) {
	companyData, err := company.adapters.GetCompanyById(req.Id)
	if err != nil {
		return nil, err
	}
	category, err := UserClient.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: int32(companyData.CategoryId),
	})
	if err != nil {
		return nil, err
	}
	profileId, err := company.adapters.GetProfileIdFromCompanyId(req.Id)
	if err != nil {
		return nil, err
	}
	address, err := company.adapters.GetAddress(profileId)
	if err != nil {
		return nil, err
	}
	addressRes := &pb.CompanyAddressResponse{
		Id:       address.ID.String(),
		Country:  address.Country,
		State:    address.State,
		District: address.District,
		City:     address.City,
	}
	links, err := company.adapters.GetAllLink(profileId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	linkRes := []*pb.CompanyLinkResponse{}
	for _, link := range links {
		linkR := &pb.CompanyLinkResponse{
			Id:    link.ID.String(),
			Title: link.Title,
			Url:   link.URL,
		}
		linkRes = append(linkRes, linkR)
	}
	jobRes := []*pb.JobResponse{}
	jobs, err := company.adapters.GetAllJobForCompany(req.Id)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
		jobR := &pb.JobResponse{
			Id:          job.JobID.String(),
			Designation: job.Designation,
			Hired:       int32(job.Hired),
			Vacancy:     int32(job.Capacity) - int32(job.Hired),
			ValidUntil:  job.ValidUntil.String(),
		}
		jobRes = append(jobRes, jobR)
	}
	res := &pb.GetCompanyResponse{
		Id:        companyData.ID.String(),
		Name:      companyData.Name,
		Phone:     companyData.Phone,
		Email:     companyData.Email,
		AvgRating: float32(companyData.AvgRating),
		Category:  category.Category,
		Address:   addressRes,
		Links:     linkRes,
		Jobs:      jobRes,
	}
	return res, nil
}
func (company *CompanyService) BlockCompany(ctx context.Context, req *pb.GetJobByCompanyId) (*emptypb.Empty, error) {
	if err := company.adapters.BlockCompany(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (company *CompanyService) UnblockCompany(ctx context.Context, req *pb.GetJobByCompanyId) (*emptypb.Empty, error) {
	if err := company.adapters.UnblockCompany(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
