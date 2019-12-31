package main

import (
	"context"
	"github.com/phamtanlong/go-crud/users/pb"
	"github.com/phamtanlong/go-crud/users/service"
)

// AuthenticationServer can be embedded to have forward compatible implementations.
type AuthenticationServer struct {
	UserService service.UserService
}

func (au *AuthenticationServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id, err := au.UserService.Register(ctx, req.Username, req.Password)

	if err != nil {
		return &pb.RegisterResponse{
			Id:    uint32(id),
			Error: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Id:    uint32(id),
		Error: "",
	}, nil
}
func (au *AuthenticationServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := au.UserService.Login(ctx, req.Username, req.Password)

	if err != nil {
		return &pb.LoginResponse{
			Token: token,
			Error: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Token: token,
		Error: "",
	}, nil
}
func (au *AuthenticationServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	id, err := au.UserService.Verify(ctx, req.GetToken())
	res := &pb.VerifyResponse{
		Id: uint32(id),
	}
	if err != nil {
		res.Error = err.Error()
	}

	return res, nil
}
