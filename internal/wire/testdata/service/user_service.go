package service

type UserService struct {
	Component           struct{}
	NotificationService *NotificationService `autowired:"true"`
}

func (s *UserService) CreateUser(username string) {
	s.NotificationService.SendNotification("Created user: " + username)
}
