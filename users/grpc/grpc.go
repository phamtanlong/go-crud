package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/phamtanlong/go-crud/users/pb"
	"github.com/phamtanlong/go-crud/users/service"
	"google.golang.org/grpc"
)

const (
	grpcPort         = "7770"
	gatewayAliveAddr = "http://localhost:8080/service/alive"
	gatewayDieAddr   = "http://localhost:8080/service/die"
)

func main() {
	// init database

	db := initDatabase()
	defer db.Close()

	// listen and serve

	userService := service.UserService{DB: db}

	// grpc
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Printf("> can not create grpc net.Listener")
	}

	authenServer := AuthenticationServer{UserService: userService}
	grpcServer := grpc.NewServer()

	pb.RegisterAuthenticationServer(grpcServer, &authenServer)

	go scheduleSendAlive(gatewayAliveAddr, "localhost", grpcPort)

	channelShutdown := make(chan os.Signal, 1)
	signal.Notify(channelShutdown, os.Interrupt)

	go func() {
		<- channelShutdown
		sendDie(gatewayDieAddr, "localhost", grpcPort)
		log.Println("Stop grpc server graceful")
		time.Sleep(2 * time.Second)
		grpcServer.Stop()
	}()

	log.Printf("> auth service start grpc %s", grpcPort)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("> can not call grpcService.Serve")
	}

}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@/go-crud?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&service.User{})

	return db
}

func sendDie(gatewayUrl string, ip string, port string) {
	err := postRequest(gatewayUrl, map[string]string{
		"kind": "user-service",
		"ip":   ip,
		"port": port,
	})

	if err != nil {
		log.Println(err)
	}
}

func scheduleSendAlive(gatewayUrl string, ip string, port string) {
	for {
		err := postRequest(gatewayUrl, map[string]string{
			"kind": "user-service",
			"ip":   ip,
			"port": port,
		})

		if err != nil {
			log.Println(err)
		}

		time.Sleep(60 * time.Second)
	}
}

func postRequest(url string, data interface{}) error {
	byteArr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteArr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	client.Do(req)
	return nil
}
