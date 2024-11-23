package service

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/cache"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/repository"
)

type EnhancedUserService struct {
	Component           struct{}
	NotificationService *NotificationService       `autowired:"true"`
	UserRepository      *repository.UserRepository `autowired:"true"`
	Logger              logger.Logger              `autowired:"true" qualifier:"console"`
	RedisCache          cache.Cache                `autowired:"true" qualifier:"redis"`
	MemoryCache         cache.Cache                `autowired:"true" qualifier:"memory"`
}

func (s *EnhancedUserService) CreateUserWithCache(username string) error {
	// Try memory cache first
	if _, err := s.MemoryCache.Get(username); err == nil {
		s.Logger.Info("User found in memory cache")
		return nil
	}

	// Try Redis cache
	if _, err := s.RedisCache.Get(username); err == nil {
		s.Logger.Info("User found in Redis cache")
		return nil
	}

	// Check repository
	if err := s.UserRepository.FindByID(username); err != nil {
		s.Logger.Error("Failed to find user in repository")
		return err
	}

	// Notify about user creation
	s.NotificationService.SendNotification("User processed: " + username)

	// Cache the result
	s.MemoryCache.Set(username, "cached")
	s.RedisCache.Set(username, "cached")

	return nil
}
