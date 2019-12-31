package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/phamtanlong/go-crud/gateway/service"
	"github.com/phamtanlong/go-crud/news"
)

type ServiceAliveRequest struct {
	Kind string `json:"kind"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
type ServiceAliveResponse struct {
	Error string `json:"error"`
}

func MakeServiceAliveEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ServiceAliveRequest)
		error := g.ServiceAlive(ctx, req.Kind, req.Ip, req.Port)
		return ServiceAliveResponse{Error: error}, nil
	}
}

type ServiceDieRequest struct {
	Kind string `json:"kind"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
type ServiceDieResponse struct {
	Error string `json:"error"`
}

func MakeServiceDieEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ServiceDieRequest)
		error := g.ServiceDie(ctx, req.Kind, req.Ip, req.Port)
		return ServiceDieResponse{Error: error}, nil
	}
}

type CreateNewsRequest struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Content   string `json:"content"`
	Tags      string `json:"tags"`
}
type CreateNewsResponse struct {
	Id    uint   `json:"id"`
	Error string `json:"error"`
}

func MakeCreateNewsEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateNewsRequest)
		id, error := g.CreateNews(ctx, req.Title, req.Thumbnail, req.Content, req.Tags)
		return CreateNewsResponse{Id: id, Error: error}, nil
	}
}

type UpdateNewsRequest struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Content   string `json:"content"`
	Tags      string `json:"tags"`
}
type UpdateNewsResponse struct {
	Error string `json:"error"`
}

func MakeUpdateNewsEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateNewsRequest)
		error := g.UpdateNews(ctx, req.Id, req.Title, req.Thumbnail, req.Content, req.Tags)
		return UpdateNewsResponse{Error: error}, nil
	}
}

type DeleteNewsRequest struct {
	Id uint `json:"id"`
}
type DeleteNewsResponse struct {
	Error string `json:"error"`
}

func MakeDeleteNewsEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteNewsRequest)
		error := g.DeleteNews(ctx, req.Id)
		return DeleteNewsResponse{Error: error}, nil
	}
}

type ReadNewsRequest struct {
	Id uint `json:"id"`
}
type ReadNewsResponse struct {
	News  *news.News `json:"news"`
	Error string     `json:"error"`
}

func MakeReadNewsEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ReadNewsRequest)
		news, error := g.ReadNews(ctx, req.Id)
		return ReadNewsResponse{News: news, Error: error}, nil
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	Id    uint   `json:"id"`
	Error string `json:"error"`
}

func MakeRegisterEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		id, error := g.Register(ctx, req.Username, req.Password)
		return RegisterResponse{Id: id, Error: error}, nil
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

func MakeLoginEndpoint(g service.GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, error := g.Login(ctx, req.Username, req.Password)
		return LoginResponse{Token: token, Error: error}, nil
	}
}

type Endpoints struct {
	ServiceAlive endpoint.Endpoint
	ServiceDie   endpoint.Endpoint
	CreateNews   endpoint.Endpoint
	UpdateNews   endpoint.Endpoint
	DeleteNews   endpoint.Endpoint
	ReadNews     endpoint.Endpoint
	Register     endpoint.Endpoint
	Login        endpoint.Endpoint
}
