package usecases

import "github.com/vishnusunil243/Job-Portal-proto-files/pb"

type Usecase interface {
	UploadImage(*pb.CompanyImageRequest, string) (string, error)
}
