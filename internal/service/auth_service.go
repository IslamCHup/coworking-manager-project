package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		firstName string,
		lastName string,
		email string,
		password string,
	) (*models.User, error)

	Login(email string, password string) (*models.User, error)
}

type authService struct {
	authRepo repository.AuthRepository
	logger   *slog.Logger
}

func NewAuthService(
	authRepo repository.AuthRepository,
	logger *slog.Logger,
) AuthService {
	return &authService{
		authRepo: authRepo,
		logger:   logger,
	}
}


func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


func (s *authService) Register(
	firstName string,
	lastName string,
	email string,
	password string,
) (*models.User, error) {

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return nil, fmt.Errorf("email: %s или пароль пустые", email)
	}

	if len(password) < 8 {
		return nil, errors.New("пароль должен быть не короче 8 символов")
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		s.logger.Error("hash password failed", "error", err)
		return nil, err
	}

	user := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		IsBlocked:    false,
		Balance:      0,
	}

	if err := s.authRepo.CreateUser(user); err != nil {
		s.logger.Error("register failed", "email", email, "error", err)
		return nil, err
	}

	s.logger.Info(
		"user registered",
		"user_id", user.ID,
		"email", email,
	)

	return user, nil
}

func (s *authService) Login(email string, password string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return nil, errors.New("email или пароль пустые")
	}

	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		s.logger.Warn("login failed - user not found", "email", email)
		return nil, errors.New("неверные учетные данные")
	}

	if user.IsBlocked {
		return nil, errors.New("пользователь заблокирован")
	}

	if err := checkPassword(user.PasswordHash, password); err != nil {
		s.logger.Warn("login failed - wrong password", "email", email)
		return nil, errors.New("неверные учетные данные")
	}

	s.logger.Info(
		"user logged in",
		"user_id", user.ID,
		"email", email,
	)

	return user, nil
}
/*
package service

import (
	"errors"

	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (uint, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (uint, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return 0, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return 0, errors.New("invalid password")
	}

	return user.ID, nil
}
*/