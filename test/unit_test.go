package companyServiceTest

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vishnusunil243/Job-Portal-Company-Service/entities"
	mock_adapters "github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters/mockAdapters"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/service"
	mock_usecases "github.com/vishnusunil243/Job-Portal-Company-Service/internal/usecases/mockUsecase"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

func TestCompanyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecase(ctrl)
	companyService := service.NewCompanyService(adapter, usecase, "", "", "")
	requestPassword := "valid"
	hashedPassword, _ := helper.HashPassword(requestPassword)
	testUUID := uuid.New()
	tests := []struct {
		name                    string
		request                 *pb.CompanyLoginRequest
		mockGetCompanyFromEmail func(string) (entities.Company, error)
		wantError               bool
		expectedResult          *pb.CompanySignupResponse
	}{
		{
			name: "Valid credentials",
			request: &pb.CompanyLoginRequest{
				Email:    "valid@gmail.com",
				Password: "valid",
			},
			mockGetCompanyFromEmail: func(s string) (entities.Company, error) {
				return entities.Company{
					ID:       testUUID,
					Name:     "valid",
					Email:    "valid@gmail.com",
					Phone:    "9999999999",
					Password: hashedPassword,
				}, nil
			},
			wantError: false,
			expectedResult: &pb.CompanySignupResponse{
				Id:    testUUID.String(),
				Email: "valid@gmail.com",
				Name:  "valid",
				Phone: "9999999999",
			},
		},
		{
			name: "Invalid Credentials",
			request: &pb.CompanyLoginRequest{
				Email:    "invalid@gmail.com",
				Password: "invalid",
			},
			mockGetCompanyFromEmail: func(s string) (entities.Company, error) {
				return entities.Company{
					Name:     "",
					Email:    "",
					Password: "",
					Phone:    "",
				}, nil
			},
			wantError: true,
			expectedResult: &pb.CompanySignupResponse{
				Id:    "id",
				Email: "asdf@gmail.com",
				Name:  "asdf",
				Phone: "888888888",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetCompanyByEmail(gomock.Any()).DoAndReturn(test.mockGetCompanyFromEmail).AnyTimes().Times(1)
			result, err := companyService.CompanyLogin(context.Background(), test.request)
			if test.wantError {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got: %v", err)
				}
				if result.Id != test.expectedResult.Id ||
					result.Email != test.expectedResult.Email ||
					result.Name != test.expectedResult.Name ||
					result.Phone != test.expectedResult.Phone {
					t.Errorf("unexpected result, got: %+v, want: %+v", result, test.expectedResult)
				}
			}
		})
	}
}
func TestCompanyAddAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecase(ctrl)
	companyservice := service.NewCompanyService(adapter, usecase, "", "", "")
	testUUid := uuid.New()
	tests := []struct {
		name                          string
		request                       *pb.CompanyAddAddressRequest
		mockGetProfileIdFromCompanyId func(string) (string, error)
		mockGetAddress                entities.Address
		wantError                     bool
	}{
		{
			name: "Success",
			request: &pb.CompanyAddAddressRequest{
				Country:   "valid",
				State:     "valid",
				District:  "valid",
				City:      "valid",
				CompanyId: testUUid.String(),
			},
			mockGetProfileIdFromCompanyId: func(s string) (string, error) {
				return testUUid.String(), nil
			},
			mockGetAddress: entities.Address{
				Country:  "",
				State:    "",
				District: "",
				City:     "",
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.CompanyAddAddressRequest{
				Country:   "invalid",
				State:     "invalid",
				District:  "invalid",
				City:      "invalid",
				CompanyId: "invalid",
			},
			mockGetProfileIdFromCompanyId: func(s string) (string, error) {
				return testUUid.String(), nil
			},
			mockGetAddress: entities.Address{
				ProfileId: testUUid,
				Country:   "valid",
				State:     "valid",
				District:  "valid",
				City:      "valid",
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetProfileIdFromCompanyId(gomock.Any()).DoAndReturn(test.mockGetProfileIdFromCompanyId).Times(1)
			adapter.EXPECT().GetAddress(gomock.Any()).Return(test.mockGetAddress, nil).Times(1)
			if !test.wantError {
				adapter.EXPECT().AddAddress(gomock.Any()).Return(nil).Times(1)
			}
			_, err := companyservice.CompanyAddAddress(context.Background(), test.request)
			if test.wantError {
				if err == nil {
					t.Errorf("expected an error but didn't get an error")
				}
			} else {
				if err != nil {
					t.Errorf("expected no errors but found %v", err)
				}
			}
		})
	}
}
func TestNotifyMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecase(ctrl)
	companyservice := service.NewCompanyService(adapter, usecase, "", "", "")
	testUUID := uuid.New()
	tests := []struct {
		name            string
		request         *pb.NotifyMeRequest
		mockGetNotifyMe func(string, string) (entities.NotifyMe, error)
		wantError       bool
	}{
		{
			name: "Fail",
			request: &pb.NotifyMeRequest{
				UserId:    testUUID.String(),
				CompanyId: testUUID.String(),
			},
			mockGetNotifyMe: func(s1, s2 string) (entities.NotifyMe, error) {
				return entities.NotifyMe{
					ID:        testUUID,
					CompanyId: testUUID,
					UserId:    testUUID,
				}, nil
			},
			wantError: true,
		},
		{
			name: "Success",
			request: &pb.NotifyMeRequest{
				UserId:    testUUID.String(),
				CompanyId: testUUID.String(),
			},
			mockGetNotifyMe: func(s1, s2 string) (entities.NotifyMe, error) {
				return entities.NotifyMe{}, nil
			},
			wantError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetNotifyMe(gomock.Any(), gomock.Any()).DoAndReturn(test.mockGetNotifyMe).AnyTimes().Times(1)
			if !test.wantError {
				adapter.EXPECT().NotifyMe(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}
			_, err := companyservice.NotifyMe(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestCompanySignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecase(ctrl)
	companyservice := service.NewCompanyService(adapter, usecase, "", "", "")
	tests := []struct {
		name                  string
		request               *pb.CompanySignupRequest
		mockGetCompanyByEmail func(string) (entities.Company, error)
		mockCompanySignup     func(entities.Company) (entities.Company, error)
		wantError             bool
		expectedResult        *pb.CompanySignupResponse
	}{
		{
			name: "Success",
			request: &pb.CompanySignupRequest{
				Email:      "valid@gmail.com",
				Name:       "valid",
				Phone:      "8888888888",
				CategoryId: 1,
				Password:   "valid",
			},
			mockGetCompanyByEmail: func(s string) (entities.Company, error) {
				return entities.Company{}, nil
			},
			mockCompanySignup: func(c entities.Company) (entities.Company, error) {
				return entities.Company{
					Email:      "valid@gmail.com",
					Name:       "valid",
					Phone:      "8888888888",
					CategoryId: 1,
					Password:   "valid",
				}, nil
			},
			wantError: false,
			expectedResult: &pb.CompanySignupResponse{
				Email:      "valid@gmail.com",
				Name:       "valid",
				Phone:      "8888888888",
				CategoryId: 1,
			},
		},
		{
			name: "Fail",
			request: &pb.CompanySignupRequest{
				Email:      "invalid@gmail.com",
				Name:       "invalid",
				Phone:      "8888888888",
				CategoryId: 1,
				Password:   "invalid",
			},
			mockGetCompanyByEmail: func(s string) (entities.Company, error) {
				return entities.Company{
					Name:       "invalid",
					Email:      "invalid@gmail.com",
					Phone:      "8888888888",
					CategoryId: 1,
				}, nil
			},
			mockCompanySignup: func(c entities.Company) (entities.Company, error) {
				return entities.Company{}, nil
			},
			wantError:      true,
			expectedResult: &pb.CompanySignupResponse{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetCompanyByEmail(gomock.Any()).DoAndReturn(test.mockGetCompanyByEmail).AnyTimes().Times(1)
			if !test.wantError {
				adapter.EXPECT().CompanySignup(gomock.Any()).DoAndReturn(test.mockCompanySignup).AnyTimes().Times(1)
			}
			res, err := companyservice.CompanySignup(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				res.Id = ""
				assert.NotNil(t, res)
				assert.Equal(t, test.expectedResult, res)
			}
		})
	}
}
