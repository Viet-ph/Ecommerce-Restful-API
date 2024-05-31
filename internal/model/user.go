package model

import (
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserCredential
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserCredential(email, password string) (*UserCredential, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &UserCredential{
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func DbUsertoUser(dbUser *db.User) User {
	return User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		UserCredential: UserCredential{
			Email:    dbUser.Email,
			Password: dbUser.Password,
		},
		Location:  dbUser.Location,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
