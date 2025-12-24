package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/auth/jwt"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
)

type RefreshService interface {
	Refresh(refreshToken string) (accessToken string, newRefreshToken string, err error)
	Logout(refreshToken string) error
	CreateForUser(userID uint) (string, error)
}

type refreshService struct {
	refreshRepo repository.RefreshTokenRepository
	logger      *slog.Logger
	refreshTTL  time.Duration
}

func NewRefreshService(
	refreshRepo repository.RefreshTokenRepository,
	logger *slog.Logger,
) RefreshService {
	return &refreshService{
		refreshRepo: refreshRepo,
		logger:      logger,
		refreshTTL:  30 * 24 * time.Hour,
	}
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *refreshService) CreateForUser(userID uint) (string, error) {
	if userID == 0 {
		return "", errors.New("invalid user id")
	}

	if err := s.refreshRepo.DeleteByUserID(userID); err != nil {
		s.logger.Error(
			"failed to delete existing refresh token",
			"user_id", userID,
			"error", err,
		)
		return "", err
	}

	rawRefresh, err := generateRefreshToken()
	if err != nil {
		return "", err
	}

	hash := hashRefreshToken(rawRefresh)

	refresh := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}

	if err := s.refreshRepo.CreateToken(refresh); err != nil {
		return "", err
	}

	s.logger.Info(
		"refresh token created",
		"user_id", userID,
	)

	return rawRefresh, nil
}

func (s *refreshService) Refresh(refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", errors.New("refresh token is empty")
	}

	hash := hashRefreshToken(refreshToken)

	stored, err := s.refreshRepo.GetByHash(hash)
	if err != nil {
		s.logger.Warn("refresh token not found", "error", err)
		return "", "", errors.New("invalid refresh token")
	}

	if time.Now().After(stored.ExpiresAt) {
		_ = s.refreshRepo.DeleteByHash(hash)
		s.logger.Warn("refresh token expired", "user_id", stored.UserID)
		return "", "", errors.New("refresh token expired")
	}

	if err := s.refreshRepo.DeleteByUserID(stored.UserID); err != nil {
		s.logger.Error(
			"failed to delete old refresh token",
			"user_id", stored.UserID,
			"error", err,
		)
		return "", "", err
	}

	rawRefresh, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	newHash := hashRefreshToken(rawRefresh)

	refresh := &models.RefreshToken{
		UserID:    stored.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}

	if err := s.refreshRepo.CreateToken(refresh); err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(stored.UserID)
	if err != nil {
		return "", "", err
	}

	s.logger.Info(
		"refresh rotation success",
		"user_id", stored.UserID,
	)

	return access, rawRefresh, nil
}

func (s *refreshService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	hash := hashRefreshToken(refreshToken)

	if err := s.refreshRepo.DeleteByHash(hash); err != nil {
		s.logger.Error("logout failed", "error", err)
		return err
	}

	s.logger.Info("logout success")
	return nil
}
