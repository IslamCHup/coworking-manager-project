package service

import (
	"errors"
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService interface {
	VerifyAdmin(login, password string) (*models.Admin, error)
	CreateAdmin(login, password string) (*models.Admin, error)
}

type adminService struct {
	repo   repository.AdminRepository
	logger *slog.Logger
}

func NewAdminService(repo repository.AdminRepository, logger *slog.Logger) AdminService {
	return &adminService{
		repo:   repo,
		logger: logger,
	}
}

func (s *adminService) VerifyAdmin(login, password string) (*models.Admin, error) {
	admin, err := s.repo.GetAdminByLogin(login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Admin not found", "login", login)
			return nil, errors.New("неверные учетные данные")
		}
		s.logger.Error("GetAdminByLogin failed", "login", login, "error", err)
		return nil, err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("Invalid password", "login", login)
		return nil, errors.New("invalid credentials")
	}

	s.logger.Info("Admin verified", "admin_id", admin.ID, "login", login)
	return admin, nil
}

func (s *adminService) CreateAdmin(login, password string) (*models.Admin, error) {
	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, err
	}

	admin := &models.Admin{
		Login:        login,
		PasswordHash: string(passwordHash),
	}

	if err := s.repo.CreateAdmin(admin); err != nil {
		s.logger.Error("CreateAdmin failed", "login", login, "error", err)
		return nil, err
	}

	s.logger.Info("Admin created", "admin_id", admin.ID, "login", login)
	return admin, nil
}
