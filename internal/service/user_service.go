package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{
		queries: q,
	}
}

func (userService *UserService) Create(ctx context.Context, loc, email, password, username string) (dto.User, error) {
	if exist, err := userService.UserExists(ctx, email); err != nil {
		return dto.User{}, fmt.Errorf("error checking if user email are already used")
	} else if exist {
		return dto.User{}, fmt.Errorf("this email is already used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return dto.User{}, fmt.Errorf("error creating new user credential")
	}

	user, err := userService.queries.CreateUser(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Location:  loc,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return dto.User{}, err
	}

	return dto.DbUsertoDto(&user), nil
}

func (userService *UserService) GetUserById(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := userService.queries.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (dto.User, error) {
	user, err := userService.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.User{}, err
	}

	return dto.DbUsertoDto(&user), nil
}

func (userService *UserService) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	err := userService.queries.DeleteUserById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash new password, %v", err)
	}

	err = userService.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       id,
		Password: string(hashedPassword),
	})
	if err != nil {
		return fmt.Errorf("unable to update new password, %v", err)
	}

	return nil
}

func (userService *UserService) UserExists(ctx context.Context, email string) (bool, error) {
	exist, err := userService.queries.UserExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("unable to check if user already created, %v", err)
	}

	return exist, nil
}
