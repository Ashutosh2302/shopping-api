package service

import (
	"database/sql"
	"errors"
	"fmt"

	"shopping_api/pkg/dto"
	"shopping_api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	db *sql.DB
}

func NewLoginService(db *sql.DB) *LoginService {
	return &LoginService{db}
}

func (s *LoginService) Login(ctx *gin.Context, u *dto.LoginRequest) (*dto.LoginResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var id string
	var password string

	query := "SELECT id, password FROM platform_user WHERE username = $1"

	err = tx.QueryRow(query, u.Username).Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if err != nil {
		fmt.Print("Invalid pasword", err)
		return nil, errors.New("invalid password")
	}

	accessToken, err := utils.GenerateAccessToken(u.Username, id)
	if err != nil {
		fmt.Print("Error generating access token:", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Id:          id,
		Username:    u.Username,
		AccessToken: accessToken,
	}, nil
}

func (s *LoginService) userExists(username string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM platform_user WHERE username = $1)"

	err := s.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to get user")
	}

	return exists, nil
}

func (s *LoginService) SignUp(ctx *gin.Context, u *dto.SignupRequest) error {
	exists, err := s.userExists(u.Username)
	if err != nil {
		return err
	}
	fmt.Print("cccc", exists)
	if exists {
		return errors.New("username already exits pleas login")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO platform_user (email, username, password) VALUES ($1, $2, $3)", u.Email, u.Username, hashedPassword)
	return err
}
