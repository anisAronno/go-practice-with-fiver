package services

import (
	"context"
	"errors"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
	"time"
)

type BlogService struct {
	blogRepo *repositories.BlogRepository
	cache    *CacheService
}

func NewBlogService(blogRepo *repositories.BlogRepository, cache *CacheService) *BlogService {
	return &BlogService{
		blogRepo: blogRepo,
		cache:    cache,
	}
}

func (s *BlogService) Create(userID uint, req *dto.CreateBlogRequest) (*models.Blog, error) {
	blog := &models.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.blogRepo.Create(blog); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.cache.InvalidateBlogCache(ctx)

	return s.blogRepo.FindByID(blog.ID)
}

func (s *BlogService) GetAll(pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetPerPage()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cacheKey := s.cache.BlogListKey(offset, limit)
	if cached, err := s.cache.GetCachedList(ctx, cacheKey); err == nil && cached != nil {

		if blogs, ok := cached.Data.([]interface{}); ok {
			result := make([]models.Blog, 0, len(blogs))
			for _, b := range blogs {
				if blogMap, ok := b.(map[string]interface{}); ok {
					blog := mapToBlog(blogMap)
					result = append(result, blog)
				}
			}
			return result, cached.Total, nil
		}
	}

	blogs, total, err := s.blogRepo.FindAll(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	blogMaps := make([]map[string]interface{}, len(blogs))
	for i, blog := range blogs {
		blogMaps[i] = blog.ToListResponse()
	}
	s.cache.SetCachedList(ctx, cacheKey, blogMaps, total)

	return blogs, total, nil
}

func mapToBlog(m map[string]interface{}) models.Blog {
	blog := models.Blog{}
	if id, ok := m["id"].(float64); ok {
		blog.ID = uint(id)
	}
	if title, ok := m["title"].(string); ok {
		blog.Title = title
	}
	if image, ok := m["image"].(string); ok {
		blog.Image = &image
	}
	if userID, ok := m["user_id"].(float64); ok {
		blog.UserID = uint(userID)
	}
	if user, ok := m["user"].(map[string]interface{}); ok {
		if uid, ok := user["id"].(float64); ok {
			blog.User.ID = uint(uid)
		}
		if name, ok := user["name"].(string); ok {
			blog.User.Name = name
		}
	}
	return blog
}

func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cacheKey := s.cache.BlogDetailKey(id)
	var cachedBlog models.Blog
	if err := s.cache.GetCachedDetail(ctx, cacheKey, &cachedBlog); err == nil && cachedBlog.ID != 0 {
		return &cachedBlog, nil
	}

	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.cache.SetCachedDetail(ctx, cacheKey, blog)

	return blog, nil
}

func (s *BlogService) GetByUserID(userID uint, pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetPerPage()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cacheKey := s.cache.UserBlogListKey(userID, offset, limit)
	if cached, err := s.cache.GetCachedList(ctx, cacheKey); err == nil && cached != nil {
		if blogs, ok := cached.Data.([]interface{}); ok {
			result := make([]models.Blog, 0, len(blogs))
			for _, b := range blogs {
				if blogMap, ok := b.(map[string]interface{}); ok {
					blog := mapToBlog(blogMap)
					result = append(result, blog)
				}
			}
			return result, cached.Total, nil
		}
	}

	blogs, total, err := s.blogRepo.FindByUserID(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	blogMaps := make([]map[string]interface{}, len(blogs))
	for i, blog := range blogs {
		blogMaps[i] = blog.ToListResponse()
	}
	s.cache.SetCachedList(ctx, cacheKey, blogMaps, total)

	return blogs, total, nil
}

func (s *BlogService) Update(id, userID uint, req *dto.UpdateBlogRequest, isAdmin bool) (*models.Blog, error) {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("blog not found")
	}

	if blog.UserID != userID && !isAdmin {
		return nil, errors.New("unauthorized")
	}

	if req.Title != "" {
		blog.Title = req.Title
	}
	if req.Content != "" {
		blog.Content = req.Content
	}

	if err := s.blogRepo.Update(blog); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.cache.InvalidateBlogDetail(ctx, id)
	s.cache.InvalidateBlogCache(ctx)

	return blog, nil
}

func (s *BlogService) Delete(id, userID uint, isAdmin bool) error {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}

	if blog.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	err = s.blogRepo.Delete(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.cache.InvalidateBlogCache(ctx)

	return nil
}

func (s *BlogService) GetAllDeleted(pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindAllDeleted(pagination.GetOffset(), pagination.GetPerPage())
}

func (s *BlogService) Restore(id uint) error {
	_, err := s.blogRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted blog not found")
	}

	err = s.blogRepo.Restore(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.cache.InvalidateBlogCache(ctx)

	return nil
}

func (s *BlogService) ForceDelete(id uint) error {
	_, err := s.blogRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted blog not found")
	}
	return s.blogRepo.ForceDelete(id)
}

func (s *BlogService) UpdateImage(id uint, imagePath string) error {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}
	blog.Image = &imagePath

	err = s.blogRepo.Update(blog)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.cache.InvalidateBlogDetail(ctx, id)
	s.cache.InvalidateBlogCache(ctx)

	return nil
}
