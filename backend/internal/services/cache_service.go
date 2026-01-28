package services

import (
	"context"
	"fmt"
	"gofiver/internal/database"
	"time"
)

type CacheService struct {
	redis   *database.RedisClient
	enabled bool
	listTTL time.Duration
}

type CachedBlogList struct {
	Blogs []CachedBlog `json:"blogs"`
	Total int64        `json:"total"`
}

type CachedBlog struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
}

func NewCacheService(redis *database.RedisClient, enabled bool) *CacheService {
	return &CacheService{
		redis:   redis,
		enabled: enabled,
		listTTL: 60 * time.Second,
	}
}

func (s *CacheService) IsEnabled() bool {
	return s.enabled
}

func (s *CacheService) GetBlogList(ctx context.Context, offset, limit int, search string) (*CachedBlogList, error) {
	if !s.enabled {
		return nil, fmt.Errorf("disabled")
	}
	key := fmt.Sprintf("bl:%d:%d:%s", offset, limit, search)
	var result CachedBlogList
	if err := s.redis.GetJSON(ctx, key, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CacheService) SetBlogList(ctx context.Context, offset, limit int, search string, data *CachedBlogList) {
	if !s.enabled {
		return
	}
	key := fmt.Sprintf("bl:%d:%d:%s", offset, limit, search)
	s.redis.SetJSON(ctx, key, data, s.listTTL)
}

func (s *CacheService) GetTotalCount(ctx context.Context) (int64, error) {
	if !s.enabled {
		return 0, fmt.Errorf("disabled")
	}
	var count int64
	if err := s.redis.GetJSON(ctx, "bl:count", &count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *CacheService) SetTotalCount(ctx context.Context, count int64) {
	if !s.enabled {
		return
	}
	s.redis.SetJSON(ctx, "bl:count", count, 5*time.Minute)
}

func (s *CacheService) InvalidateAll(ctx context.Context) {
	if !s.enabled {
		return
	}
	s.redis.DeletePattern(ctx, "bl:*")
}
