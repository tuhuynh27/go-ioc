package repository

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/config"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
)

type UserRepository struct {
	Component          struct{}
	Logger             logger.Logger              `autowired:"true" qualifier:"json"`
	Config             *config.DatabaseConfigData `autowired:"true"`
	DatabaseConfigData *config.DatabaseConfigData `autowired:"true"`
}

func (r *UserRepository) FindByID(id string) error {
	r.Logger.Info("Finding user with ID: " + id + " using DSN: " + r.DatabaseConfigData.GetConfig().GetDSN())
	return nil
}
