package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/cache"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/metrics"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	ioc.Component
	Logger      logger.Logger            `autowired:"" qualifier:"json"`
	Cache       cache.Cache              `autowired:""`
	Metrics     metrics.MetricsCollector `autowired:""`
	EmailSender MessageService           `autowired:"" qualifier:"email"`
	SmsSender   MessageService           `autowired:"" qualifier:"sms"`
}

func (s *UserService) CreateUser(user User) error {
	// Log the operation
	userJson, _ := json.Marshal(user)
	s.Logger.LogWithLevel(logger.INFO, fmt.Sprintf("Creating user: %s", string(userJson)))

	// Cache the user
	s.Cache.Set("user:"+user.ID, user, 24*time.Hour)

	// Record metrics
	s.Metrics.RecordMetric("users.created", 1)

	// Send welcome messages
	go s.EmailSender.Send(user.Email, "Welcome to our service!")
	go s.SmsSender.Send(user.Email, "Welcome! Your account has been created.")

	return nil
}

func (s *UserService) GetUser(id string) (*User, error) {
	// Try to get from cache
	if cached, found := s.Cache.Get("user:" + id); found {
		if user, ok := cached.(User); ok {
			s.Logger.Log("User found in cache: " + id)
			s.Metrics.RecordMetric("cache.hits", 1)
			return &user, nil
		}
	}

	s.Logger.LogWithLevel(logger.ERROR, "User not found: "+id)
	s.Metrics.RecordMetric("cache.misses", 1)
	return nil, fmt.Errorf("user not found")
}
