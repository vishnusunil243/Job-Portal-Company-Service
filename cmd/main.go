package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-Company-Service/db"
	"github.com/vishnusunil243/Job-Portal-Company-Service/initializer"
	"github.com/vishnusunil243/Job-Portal-Company-Service/internal/service"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatalf(err.Error())
	}
	userConn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	searchConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while connecting to user service")
	}
	defer func() {
		userConn.Close()
		searchConn.Close()
	}()
	userRes := pb.NewUserServiceClient(userConn)
	searchRes := pb.NewSearchServiceClient(searchConn)
	service.UserClient = userRes
	service.SearchClient = searchRes
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterCompanyServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen on port 8082")
	}
	log.Printf("company service listening on port 8082")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to listen on port 8082")
	}
}
