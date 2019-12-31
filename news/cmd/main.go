package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/phamtanlong/go-crud/news/endpoints"
	http1 "github.com/phamtanlong/go-crud/news/http"
	service "github.com/phamtanlong/go-crud/news/service"
	"github.com/phamtanlong/go-crud/users/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	http "net/http"
)

const (
	httpAddr = ":8888"
	grpcAddr = ":8880"
	authenticationGrpcAddr = ":7770"
)

func main() {

	// init database

	db := initDatabase()
	defer db.Close()

	// grpc authentication client
	grpcAuthConn, authenClient := initAuthenticationClient()
	defer grpcAuthConn.Close()

	// listen and serve

	newService := service.NewsService{DB:db, AuthenClient:authenClient}

	points := endpoints.Endpoints{
		Create: endpoints.MakeCreateEndpoint(newService),
		Update: endpoints.MakeUpdateEndpoint(newService),
		Delete: endpoints.MakeDeleteEndpoint(newService),
		Read:   endpoints.MakeReadEndpoint(newService),
	}

	// attach endpoint middleware here
	//like: points.Create = middleware(points.Create)

	handler := http1.NewHTTPHandler(points)

	listener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		panic("> can not create net.Listener")
	}

	log.Printf("> news service start http %s", httpAddr)
	http.Serve(listener, handler)
}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@/go-crud?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&service.News{})

	return db
}

func initAuthenticationClient() (*grpc.ClientConn, pb.AuthenticationClient) {
	conn, err := grpc.Dial(authenticationGrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return conn, nil
	}
	client := pb.NewAuthenticationClient(conn)
	return conn, client
}


