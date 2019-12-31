package news

import (
	"context"
	"github.com/jinzhu/gorm"
)

// NewsService is a service to serve news
type NewsService interface {
	Create(ctx context.Context, userId uint, title string, thumbnail string, content string, tags string) (id uint, error error)
	Update(ctx context.Context, userId uint, id uint, title string, thumbnail string, content string, tags string) (error string)
	Delete(ctx context.Context, userId uint, id uint) (error string)
	Read(ctx context.Context, userId uint, id uint) (news *News, error string)
}

// News is base model for this module
type News struct {
	gorm.Model
	Title     string
	Thumbnail string
	Content   string
	Tags      string
	CreatedBy uint
	UpdatedBy uint
}

