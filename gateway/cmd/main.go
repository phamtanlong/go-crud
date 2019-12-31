package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/phamtanlong/go-crud/gateway/endpoints"
	http1 "github.com/phamtanlong/go-crud/gateway/http"
	"github.com/phamtanlong/go-crud/gateway/service"
	"github.com/phamtanlong/go-crud/users/pb"
	"log"
	"net"
	"net/http"
)

const (
	httpPort = "8080"
)

func main() {

	db := initDatabase()
	defer db.Close()

	gatewayService := service.GatewayService{
		DB: db,
	}
	defer gatewayService.OnDestroy()

	points := endpoints.Endpoints{
		ServiceAlive: endpoints.MakeServiceAliveEndpoint(gatewayService),
		ServiceDie:   endpoints.MakeServiceDieEndpoint(gatewayService),
		CreateNews:   endpoints.MakeCreateNewsEndpoint(gatewayService),
		UpdateNews:   endpoints.MakeUpdateNewsEndpoint(gatewayService),
		DeleteNews:   endpoints.MakeDeleteNewsEndpoint(gatewayService),
		ReadNews:     endpoints.MakeReadNewsEndpoint(gatewayService),
		Register:     endpoints.MakeRegisterEndpoint(gatewayService),
		Login:        endpoints.MakeLoginEndpoint(gatewayService),
	}

	// attach middleware here
	points.CreateNews = authMiddleware(gatewayService)(points.CreateNews)
	points.UpdateNews = authMiddleware(gatewayService)(points.UpdateNews)
	points.DeleteNews = authMiddleware(gatewayService)(points.DeleteNews)

	listener, err := net.Listen("tcp", ":" + httpPort)
	if err != nil {
		log.Fatalf("> can not listen to %s", httpPort)
	}

	handler := http1.NewHTTPHandler(points)

	log.Printf("> gateway start http %s", httpPort)
	http.Serve(listener, handler)

}

func authMiddleware(gatewayService service.GatewayService) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			gatewayService.InitIfNeed()

			//get token from header
			token, ok := ctx.Value(httptransport.ContextKeyRequestAuthorization).(string)
			if !ok {
				return nil, errors.New("token not found")
			}

			//grpc to auth service
			res, err := gatewayService.AuthClient.Verify(ctx, &pb.VerifyRequest{Token: token})
			if err != nil {
				return 0, errors.New("can not connect auth service")
			}
			if len(res.Error) > 0 {
				return 0, errors.New(res.Error)
			}

			//attach user id to context
			ctx = context.WithValue(ctx, "UserId", uint(res.Id))

			return next(ctx, request)
		}
	}
}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@/go-crud?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatalln(err)
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&service.ServiceLocation{})

	return db
}
