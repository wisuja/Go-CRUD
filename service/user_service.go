package service

import (
	"context"

	"github.com/wisuja/crud/database"
	"github.com/wisuja/crud/entity"
	"github.com/wisuja/crud/repository"
)

func FetchAllUsers() ([]entity.User, error) {
	var users []entity.User

	credentials, err := database.GetDefaultDatabaseConfig("../.env")

	if err != nil {
		return users, err
	}

	db, err := database.GetConnection(credentials)

	if err != nil {
		return users, err
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	ctx := context.Background()

	users, err = userRepository.FindAll(ctx)
	if err != nil {
		return users, err
	}

	return users, nil
}

func FetchUser(id int) (entity.User, error) {
	var user entity.User

	credentials, err := database.GetDefaultDatabaseConfig("../.env")

	if err != nil {
		return user, err
	}

	db, err := database.GetConnection(credentials)

	if err != nil {
		return user, err
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	ctx := context.Background()

	user, err = userRepository.FindById(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(username, password string) (entity.User, error) {
	user := entity.User{
		Username: username,
		Password: password,
	}

	credentials, err := database.GetDefaultDatabaseConfig("../.env")

	if err != nil {
		return user, err
	}

	db, err := database.GetConnection(credentials)

	if err != nil {
		return user, err
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	ctx := context.Background()

	user, err = userRepository.Insert(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(id int, username, password string) (entity.User, error) {
	user := entity.User{
		Username: username,
		Password: password,
	}

	credentials, err := database.GetDefaultDatabaseConfig("../.env")

	if err != nil {
		return user, err
	}

	db, err := database.GetConnection(credentials)

	if err != nil {
		return user, err
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	ctx := context.Background()

	user, err = userRepository.Update(ctx, id, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func DeleteUser(id int) (bool, error) {
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

	success, err := userRepository.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return success, nil
}
