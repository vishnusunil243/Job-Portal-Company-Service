package usecases

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type CompanyUsecase struct {
	adapters adapters.AdapterInterface
}

func NewCompanyUseCase(adapters adapters.AdapterInterface) *CompanyUsecase {
	return &CompanyUsecase{
		adapters: adapters,
	}
}
func (company *CompanyUsecase) UploadImage(req *pb.CompanyImageRequest, profileId string) (string, error) {
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESSKEY"), os.Getenv("MINIO_SECRETKEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Print("error while initialising minio", err)
		return "", err
	}
	objectName := "images/" + req.ObjectName
	contentType := `image/jpeg`
	n, err := minioClient.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, bytes.NewReader(req.ImageData), int64(len(req.ImageData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error while uploading to minio", err)
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %v\n", objectName, n)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Println("error while generating presigned URL", err)
		return "", err
	}
	url, err := company.adapters.UploadImage(presignedURL.String(), profileId)
	return url, err
}
