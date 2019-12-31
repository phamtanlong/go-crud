package main

import (
	"context"
	"errors"
	"github.com/phamtanlong/go-crud/news/pb"
	"github.com/phamtanlong/go-crud/news/service"
)

type NewsServer struct {
	NewsService service.NewsService
}

func (n *NewsServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := n.NewsService.Create(ctx, uint(req.UserId), req.Title, req.Thumbnail, req.Content, req.Tags)
	if len(err) > 0 {
		return nil, errors.New(err)
	}

	return &pb.CreateResponse{
		Id:    uint32(id),
		Error: err,
	}, nil
}
func (n *NewsServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	err := n.NewsService.Update(ctx, uint(req.UserId), uint(req.Id), req.Title, req.Thumbnail, req.Content, req.Tags)
	return &pb.UpdateResponse{
		Error: err,
	}, nil
}
func (n *NewsServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := n.NewsService.Delete(ctx, uint(req.UserId), uint(req.Id))
	return &pb.DeleteResponse{
		Error: err,
	}, nil
}
func (n *NewsServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	news, err := n.NewsService.Read(ctx, 0, uint(req.Id))
	if len(err) > 0 {
		return nil, errors.New(err)
	}

	return &pb.ReadResponse{
		Id:        uint32(news.ID),
		Title:     news.Title,
		Thumbnail: news.Thumbnail,
		Content:   news.Content,
		Tags:      news.Tags,
	}, nil
}
