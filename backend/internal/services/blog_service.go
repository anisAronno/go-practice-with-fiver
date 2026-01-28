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
	repo  *repositories.BlogRepository
	cache *CacheService
}

func NewBlogService(repo *repositories.BlogRepository, cache *CacheService) *BlogService {
	return &BlogService{repo: repo, cache: cache}
}

func (s *BlogService) Create(userID uint, req *dto.CreateBlogRequest) (*models.Blog, error) {
	blog := &models.Blog{Title: req.Title, Content: req.Content, UserID: userID}
	if err := s.repo.Create(blog); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.cache.InvalidateAll(ctx)
	return s.repo.FindByID(blog.ID)
}

func (s *BlogService) GetAll(q *dto.BlogSearchQuery) ([]CachedBlog, int64, error) {
	offset := q.GetOffset()
	limit := q.GetPerPage()
	search := q.Search

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if cached, err := s.cache.GetBlogList(ctx, offset, limit, search); err == nil {
		return cached.Blogs, cached.Total, nil
	}

	blogs, err := s.repo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if search == "" {
		if cachedCount, err := s.cache.GetTotalCount(ctx); err == nil {
			total = cachedCount
		} else {
			total, _ = s.repo.GetCount("")
			s.cache.SetTotalCount(ctx, total)
		}
	} else {
		if len(blogs) < limit {
			total = int64(offset + len(blogs))
		} else {
			total = int64(offset + limit + 1)
		}
	}

	result := make([]CachedBlog, len(blogs))
	for i, b := range blogs {
		result[i] = CachedBlog{
			ID:        b.ID,
			Title:     b.Title,
			Image:     ptrToStr(b.Image),
			UserID:    b.UserID,
			UserName:  b.User.Name,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		}
	}

	s.cache.SetBlogList(ctx, offset, limit, search, &CachedBlogList{Blogs: result, Total: total})
	return result, total, nil
}

func ptrToStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	return s.repo.FindByID(id)
}

func (s *BlogService) GetByUserID(userID uint, q *dto.PaginationQuery) ([]models.Blog, int64, error) {
	blogs, err := s.repo.FindByUserID(userID, q.GetOffset(), q.GetPerPage())
	if err != nil {
		return nil, 0, err
	}
	total, _ := s.repo.GetUserBlogCount(userID)
	return blogs, total, nil
}

func (s *BlogService) Update(id, userID uint, req *dto.UpdateBlogRequest, isAdmin bool) (*models.Blog, error) {
	blog, err := s.repo.FindByID(id)
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
	if err := s.repo.Update(blog); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.cache.InvalidateAll(ctx)
	return blog, nil
}

func (s *BlogService) Delete(id, userID uint, isAdmin bool) error {
	blog, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}
	if blog.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.cache.InvalidateAll(ctx)
	return nil
}

func (s *BlogService) GetAllDeleted(q *dto.PaginationQuery) ([]models.Blog, int64, error) {
	blogs, err := s.repo.FindAllDeleted(q.GetOffset(), q.GetPerPage())
	if err != nil {
		return nil, 0, err
	}
	total, _ := s.repo.GetDeletedCount()
	return blogs, total, nil
}

func (s *BlogService) Restore(id uint) error {
	if _, err := s.repo.FindDeletedByID(id); err != nil {
		return errors.New("deleted blog not found")
	}
	if err := s.repo.Restore(id); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.cache.InvalidateAll(ctx)
	return nil
}

func (s *BlogService) ForceDelete(id uint) error {
	if _, err := s.repo.FindDeletedByID(id); err != nil {
		return errors.New("deleted blog not found")
	}
	return s.repo.ForceDelete(id)
}

func (s *BlogService) UpdateImage(id uint, imagePath string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("blog not found")
	}
	if err := s.repo.UpdateImage(id, imagePath); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.cache.InvalidateAll(ctx)
	return nil
}
