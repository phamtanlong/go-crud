package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/phamtanlong/go-crud/users/pb"
	"log"
	"time"
)

type News struct {
	gorm.Model
	Title     string
	Thumbnail string
	Content   string
	Tags      string
	CreatedBy *uint
	UpdatedBy *uint
}
type NewsService struct {
	DB *gorm.DB
	AuthenClient pb.AuthenticationClient
}

func (n NewsService) Create(ctx context.Context, userId uint, title string, thumbnail string, content string, tags string) (id uint, error string) {
	log.Printf("> create \ntitle: %s\nthumbnail: %s\ncontent: %s\ntags: %s", title, thumbnail, content, tags)

	record := News{
		Title:     title,
		Thumbnail: thumbnail,
		Content:   content,
		Tags:      tags,
		CreatedBy: &userId,
	}

	n.DB.Create(&record)
	return record.ID, ""
}
func (n NewsService) Update(ctx context.Context, userId uint, id uint, title string, thumbnail string, content string, tags string) (error string) {
	log.Printf("> update id %d\ntitle: %s\nthumbnail: %s\ncontent: %s\ntags: %s", id, title, thumbnail, content, tags)

	var record News
	if n.DB.First(&record, id).RecordNotFound() {
		return "record not found"
	}

	record.Title = title
	record.Thumbnail = thumbnail
	record.Content = content
	record.Tags = tags
	record.UpdatedBy = &userId

	n.DB.Save(&record)
	return ""
}
func (n NewsService) Delete(ctx context.Context, userId uint, id uint) (error string) {
	log.Printf("> delete id %d", id)

	var record News
	if n.DB.First(&record, id).RecordNotFound() {
		return "record not found"
	}

	record.UpdatedBy = &userId

	now := time.Now()
	record.DeletedAt = &now
	n.DB.Save(&record)

	return ""
}
func (n NewsService) Read(ctx context.Context, userId uint, id uint) (news *News, error string) {
	log.Printf("> read id %d", id)

	var record News
	if n.DB.First(&record, id).RecordNotFound() {
		return nil, "record not found"
	}
	return &record, ""
}
