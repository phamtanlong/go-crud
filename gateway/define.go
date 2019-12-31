package gateway

import (
	"context"
	"github.com/phamtanlong/go-crud/news"
)

type GatewayService interface {

	// service discovery
	//heartbeat and register also
	//TODO: need authentication on these 2 apis
	ServiceAlive(ctx context.Context, kind string, ip string, port string) (error string)
	ServiceDie(ctx context.Context, kind string, ip string, port string) (error string)


	// news service
	CreateNews(ctx context.Context, title string, thumbnail string, content string, tags string) (id uint, error string)
	UpdateNews(ctx context.Context, id uint, title string, thumbnail string, content string, tags string) (error string)
	DeleteNews(ctx context.Context, id uint) (error string)
	ReadNews(ctx context.Context, id uint) (news *news.News, error string)

	// auth service
	Register(ctx context.Context, username string, password string) (id uint, error string)
	Login(ctx context.Context, username string, password string) (token string, error string)

}
