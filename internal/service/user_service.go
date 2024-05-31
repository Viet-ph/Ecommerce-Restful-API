package service

import (
	"context"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/model"
	"github.com/google/uuid"
)

type UserService struct {
	Queries *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{
		Queries: q,
	}
}

func (userService *UserService) Create(
	ctx context.Context,
	credential model.UserCredential,
	loc, username string,
) (model.User, error) {
	user, err := userService.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Username:  username,
		Email:     credential.Email,
		Password:  credential.Password,
		Location:  loc,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return model.User{}, err
	}

	return model.DbUsertoUser(&user), nil
}

func (userService *UserService) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := userService.Queries.GetUserById(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return model.DbUsertoUser(&user), nil
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := userService.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	return model.DbUsertoUser(&user), nil
}

func (userService *UserService) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	err := userService.Queries.DeleteUserById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	err := userService.Queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:        id,
		Password:  newPassword,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) UserExists(ctx context.Context, email string) (bool, error) {
	exist, err := userService.Queries.UserExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exist, nil
}
