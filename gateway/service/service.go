package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/phamtanlong/go-crud/news"
	pb1 "github.com/phamtanlong/go-crud/news/pb"
	"github.com/phamtanlong/go-crud/users/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ServiceLocation struct {
	gorm.Model
	Kind     string    `json:"kind"`
	Ip       string    `json:"ip"`
	Port     string    `json:"port"`
	LastBeat time.Time `json:"last_beat"`
}

type GatewayService struct {
	DB *gorm.DB

	AuthConn   *grpc.ClientConn
	AuthClient pb.AuthenticationClient

	NewsConn   *grpc.ClientConn
	NewsClient pb1.NewsClient
}

func (g GatewayService) OnDestroy() {
	g.AuthConn.Close()
	g.NewsConn.Close()
}

func CreateAuthClient(DB *gorm.DB) (*grpc.ClientConn, pb.AuthenticationClient) {
	var locations = FindServiceLocation(DB, "user-service")
	if len(locations) == 0 {
		return nil, nil
	}
	var addr = locations[0].Ip + ":" + locations[0].Port

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Println(err)
		return nil, nil
	}
	client := pb.NewAuthenticationClient(conn)
	return conn, client
}

func CreateNewsClient(DB *gorm.DB) (*grpc.ClientConn, pb1.NewsClient) {
	var locations = FindServiceLocation(DB, "news-service")
	if len(locations) == 0 {
		return nil, nil
	}
	var addr = locations[0].Ip + ":" + locations[0].Port

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Println(err)
		return nil, nil
	}
	client := pb1.NewNewsClient(conn)
	return conn, client
}

func FindServiceLocation(DB *gorm.DB, kind string) []ServiceLocation {
	var locations []ServiceLocation
	if DB.Where("kind = ?", kind).Find(&locations).RecordNotFound() {
		return make([]ServiceLocation, 0)
	}
	return locations
}

func (g *GatewayService) InitIfNeed() {
	if g.NewsClient == nil {
		g.NewsConn, g.NewsClient = CreateNewsClient(g.DB)
		//TODO: create dummy client ??
	}

	if g.AuthClient == nil {
		g.AuthConn, g.AuthClient = CreateAuthClient(g.DB)
		//TODO: create dummy client ??
	}
}

func (g GatewayService) ServiceAlive(ctx context.Context, kind string, ip string, port string) string {
	log.Println("> service alive:", kind, ip, port)
	var location ServiceLocation
	var query = g.DB.Where("kind = ? AND ip = ? AND port = ?", kind, ip, port)
	if query.First(&location).RecordNotFound() {
		// add new service
		location = ServiceLocation{
			Kind:     kind,
			Ip:       ip,
			Port:     port,
			LastBeat: time.Now(),
		}
		g.DB.Save(&location)
	} else {
		// update last beat
		location.LastBeat = time.Now()
		g.DB.Save(&location)
	}
	return ""
}
func (g GatewayService) ServiceDie(ctx context.Context, kind string, ip string, port string) string {
	log.Println("> service die:", kind, ip, port)
	var location ServiceLocation
	var query = g.DB.Where("kind = ? AND ip = ? AND port = ?", kind, ip, port)
	if query.First(&location).RecordNotFound() {
		return "service not found"
	} else {
		g.DB.Delete(&location)
		return ""
	}
}
func (g GatewayService) CreateNews(ctx context.Context, title string, thumbnail string, content string, tags string) (uint, string) {
	g.InitIfNeed()
	userId, _ := ctx.Value("UserId").(uint)

	res, err := g.NewsClient.Create(ctx, &pb1.CreateRequest{
		UserId:    uint32(userId),
		Title:     title,
		Thumbnail: thumbnail,
		Content:   content,
		Tags:      tags,
	})

	if err != nil {
		return 0, err.Error()
	}

	return uint(res.Id), res.Error
}
func (g GatewayService) UpdateNews(ctx context.Context, id uint, title string, thumbnail string, content string, tags string) string {
	g.InitIfNeed()
	userId, _ := ctx.Value("UserId").(uint)
	res, err := g.NewsClient.Update(ctx, &pb1.UpdateRequest{
		Id:        uint32(id),
		UserId:    uint32(userId),
		Title:     title,
		Thumbnail: thumbnail,
		Content:   content,
		Tags:      tags,
	})

	if err != nil {
		return err.Error()
	}

	return res.Error
}
func (g GatewayService) DeleteNews(ctx context.Context, id uint) string {
	g.InitIfNeed()
	userId, _ := ctx.Value("UserId").(uint)
	res, err := g.NewsClient.Delete(ctx, &pb1.DeleteRequest{
		Id:     uint32(id),
		UserId: uint32(userId),
	})

	if err != nil {
		return err.Error()
	}

	return res.Error
}
func (g GatewayService) ReadNews(ctx context.Context, id uint) (*news.News, string) {
	g.InitIfNeed()
	res, err := g.NewsClient.Read(ctx, &pb1.ReadRequest{
		Id: uint32(id),
	})

	if err != nil {
		log.Printf("=> %s", err.Error())
		return nil, err.Error()
	}

	return &news.News{
		Title:     res.Title,
		Thumbnail: res.Thumbnail,
		Content:   res.Content,
		Tags:      res.Tags,
	}, ""
}
func (g GatewayService) Register(ctx context.Context, username string, password string) (uint, string) {
	g.InitIfNeed()
	res, err := g.AuthClient.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		return 0, err.Error()
	}

	return uint(res.Id), res.Error
}
func (g GatewayService) Login(ctx context.Context, username string, password string) (string, string) {
	g.InitIfNeed()
	res, err := g.AuthClient.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		return "", err.Error()
	}

	return res.Token, res.Error
}
