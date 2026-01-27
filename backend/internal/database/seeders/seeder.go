package seeders

import (
	"gofiver/internal/models"
	"log"

	"gorm.io/gorm"
)

// Seeder handles database seeding (like Laravel db:seed)
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Run executes all seeders
func (s *Seeder) Run() error {
	log.Println("Running database seeders...")

	// Check if already seeded
	var count int64
	s.db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// Run seeders
	if err := s.seedUsers(); err != nil {
		return err
	}

	if err := s.seedBlogs(); err != nil {
		return err
	}

	log.Println("Database seeding completed!")
	return nil
}

// seedUsers creates demo users
func (s *Seeder) seedUsers() error {
	users := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"Admin User", "admin@example.com", "password123"},
		{"John Doe", "john@example.com", "password123"},
		{"Jane Smith", "jane@example.com", "password123"},
	}

	for _, u := range users {
		user := &models.User{
			Name:  u.Name,
			Email: u.Email,
		}
		if err := user.HashPassword(u.Password); err != nil {
			return err
		}
		if err := s.db.Create(user).Error; err != nil {
			return err
		}
		log.Printf("Created user: %s (%s)", user.Name, user.Email)
	}

	return nil
}

// seedBlogs creates demo blogs
func (s *Seeder) seedBlogs() error {
	blogs := []struct {
		Title   string
		Content string
		UserID  uint
	}{
		{
			Title:   "Getting Started with Go",
			Content: "Go is a statically typed, compiled language designed at Google. It is syntactically similar to C, but with memory safety and garbage collection.",
			UserID:  1,
		},
		{
			Title:   "Understanding Fiber Framework",
			Content: "Fiber is an Express inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go. Designed to ease things up for fast development.",
			UserID:  1,
		},
		{
			Title:   "Building REST APIs in Go",
			Content: "Learn how to build production-ready REST APIs using Go and Fiber framework. This guide covers authentication, database integration, and best practices.",
			UserID:  2,
		},
		{
			Title:   "Docker for Go Applications",
			Content: "Containerizing Go applications with Docker makes deployment consistent and scalable. Learn the best practices for multi-stage builds and optimization.",
			UserID:  2,
		},
		{
			Title:   "MySQL with GORM",
			Content: "GORM is the fantastic ORM library for Go. It aims to be developer friendly with features like Auto Migrations, Relations, and Hooks.",
			UserID:  3,
		},
	}

	for _, b := range blogs {
		blog := &models.Blog{
			Title:   b.Title,
			Content: b.Content,
			UserID:  b.UserID,
		}
		if err := s.db.Create(blog).Error; err != nil {
			return err
		}
		log.Printf("Created blog: %s", blog.Title)
	}

	return nil
}
