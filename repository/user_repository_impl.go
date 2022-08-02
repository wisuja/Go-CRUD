package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/wisuja/crud/entity"
)

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (repository *userRepositoryImpl) FindAll(ctx context.Context) ([]entity.User, error) {
	query := "SELECT id, username FROM users"
	rows, err := repository.DB.QueryContext(ctx, query)

	var users []entity.User

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.Id, &user.Username)

		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository *userRepositoryImpl) FindById(ctx context.Context, id int) (entity.User, error) {
	query := "SELECT id, username FROM users WHERE id = ? LIMIT 1"
	rows, err := repository.DB.QueryContext(ctx, query, id)

	if err != nil {
		return entity.User{}, err
	}

	defer rows.Close()

	user := entity.User{}
	if !rows.Next() {
		return user, errors.New("No user found!")
	}

	err = rows.Scan(&user.Id, &user.Username)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepositoryImpl) FindByUser(ctx context.Context, user entity.User) (entity.User, error) {
	query := "SELECT id, username FROM users WHERE username = ? AND password = ? LIMIT 1"
	rows, err := repository.DB.QueryContext(ctx, query, user.Username, user.Password)

	if err != nil {
		return entity.User{}, err
	}

	defer rows.Close()

	result := entity.User{}
	if !rows.Next() {
		return result, errors.New("No user found!")
	}

	err = rows.Scan(&result.Id, &result.Username)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (repository *userRepositoryImpl) Insert(ctx context.Context, user entity.User) (entity.User, error) {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"

	result, err := repository.DB.ExecContext(ctx, query, user.Username, user.Password)

	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return user, err
	}

	user.Id = int32(id)

	return user, nil
}

func (repository *userRepositoryImpl) Update(ctx context.Context, id int, user entity.User) (entity.User, error) {
	query := "UPDATE users SET username = ?, password = ? WHERE id = ?"

	result, err := repository.DB.ExecContext(ctx, query, user.Username, user.Password, id)

	if err != nil {
		return user, err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return user, err
	}

	if rows <= 0 {
		return user, errors.New("No change")
	}

	return user, nil
}

func (repository *userRepositoryImpl) Delete(ctx context.Context, id int) (bool, error) {
	query := "DELETE FROM users WHERE id = ?"

	result, err := repository.DB.ExecContext(ctx, query, id)

	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	if rows <= 0 {
		return false, errors.New("No change")
	}

	return true, nil
}
