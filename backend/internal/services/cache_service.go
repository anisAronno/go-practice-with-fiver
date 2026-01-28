package services

import (
	"context"
	"fmt"
	"gofiver/internal/config"
	"gofiver/internal/database"
	"time"
)

type CacheService struct {
	redis  *database.RedisClient
	config *config.CacheConfig
}

type CachedList struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

func NewCacheService(redis *database.RedisClient, cfg *config.CacheConfig) *CacheService {
	return &CacheService{
		redis:  redis,
		config: cfg,
	}
}

func (s *CacheService) IsEnabled() bool {
	return s.config.Enabled
}

func (s *CacheService) BlogListKey(offset, limit int) string {
	return fmt.Sprintf("blogs:list:%d:%d", offset, limit)
}

func (s *CacheService) BlogCountKey() string {
	return "blogs:count"
}

func (s *CacheService) BlogDetailKey(id uint) string {
	return fmt.Sprintf("blogs:detail:%d", id)
}

func (s *CacheService) UserBlogListKey(userID uint, offset, limit int) string {
	return fmt.Sprintf("blogs:user:%d:%d:%d", userID, offset, limit)
}

func (s *CacheService) UserBlogCountKey(userID uint) string {
	return fmt.Sprintf("blogs:user:%d:count", userID)
}

func (s *CacheService) GetCachedList(ctx context.Context, key string) (*CachedList, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("cache disabled")
	}

	var result CachedList
	err := s.redis.GetJSON(ctx, key, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CacheService) SetCachedList(ctx context.Context, key string, data interface{}, total int64) error {
	if !s.config.Enabled {
		return nil
	}

	cached := CachedList{
		Data:  data,
		Total: total,
	}
	return s.redis.SetJSON(ctx, key, cached, s.config.ListTTL)
}

func (s *CacheService) GetCachedCount(ctx context.Context, key string) (int64, error) {
	if !s.config.Enabled {
		return 0, fmt.Errorf("cache disabled")
	}

	var count int64
	err := s.redis.GetJSON(ctx, key, &count)
	return count, err
}

func (s *CacheService) SetCachedCount(ctx context.Context, key string, count int64) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.SetJSON(ctx, key, count, s.config.CountTTL)
}

func (s *CacheService) GetCachedDetail(ctx context.Context, key string, dest interface{}) error {
	if !s.config.Enabled {
		return fmt.Errorf("cache disabled")
	}
	return s.redis.GetJSON(ctx, key, dest)
}

func (s *CacheService) SetCachedDetail(ctx context.Context, key string, data interface{}) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.SetJSON(ctx, key, data, s.config.DetailTTL)
}

func (s *CacheService) InvalidateBlogCache(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.DeletePattern(ctx, "blogs:*")
}

func (s *CacheService) InvalidateBlogDetail(ctx context.Context, id uint) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.Delete(ctx, s.BlogDetailKey(id))
}

func (s *CacheService) InvalidateUserBlogCache(ctx context.Context, userID uint) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.DeletePattern(ctx, fmt.Sprintf("blogs:user:%d:*", userID))
}

func (s *CacheService) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !s.config.Enabled {
		return nil
	}
	return s.redis.SetJSON(ctx, key, value, ttl)
}

func (s *CacheService) GetConfig() *config.CacheConfig {
	return s.config
}
