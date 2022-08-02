package service

import (
	"context"

	"github.com/wisuja/crud/database"
	"github.com/wisuja/crud/entity"
	"github.com/wisuja/crud/repository"
)

func CheckLogin(username, password string) (bool, error) {
	credentials, err := database.GetDefaultDatabaseConfig("../.env")

	if err != nil {
		return false, err
	}

	db, err := database.GetConnection(credentials)

	if err != nil {
		return false, err
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	ctx := context.Background()

	_, err = userRepository.FindByUser(ctx, entity.User{
		Username: username,
		Password: password,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
