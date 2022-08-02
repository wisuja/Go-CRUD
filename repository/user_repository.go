package repository

import (
	"context"

	"github.com/wisuja/crud/entity"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	FindById(ctx context.Context, id int) (entity.User, error)
	FindByUser(ctx context.Context, user entity.User) (entity.User, error)
	Insert(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, id int, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id int) (bool, error)
}
