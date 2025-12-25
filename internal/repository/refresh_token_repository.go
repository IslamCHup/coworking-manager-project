package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	CreateToken(token *models.RefreshToken) error
	GetByHash(hash string) (*models.RefreshToken, error)
	DeleteByUserID(userID uint) error
	DeleteByHash(hash string) error
}

type refreshTokenRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRefreshTokenRepository(db *gorm.DB, logger *slog.Logger) RefreshTokenRepository {
	return &refreshTokenRepository{
		db:     db,
		logger: logger,
	}
}

func (r *refreshTokenRepository) CreateToken(token *models.RefreshToken) error {
	if err := r.db.Create(token).Error; err != nil {
		r.logger.Error(
			"create refresh token failed",
			"user_id", token.UserID,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"refresh token created",
		"user_id", token.UserID,
	)
	return nil
}

func (r *refreshTokenRepository) GetByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	if err := r.db.Where("token_hash = ?", hash).First(&token).Error; err != nil {
		r.logger.Warn("refresh token not found", "hash", hash, "error", err)
		return nil, err
	}

	return &token, nil
}

func (r *refreshTokenRepository) DeleteByHash(hash string) error {
	err := r.db.Unscoped().
		Where("token_hash = ?", hash).
		Delete(&models.RefreshToken{}).Error
	if err != nil {
		r.logger.Error("delete refresh token by hash failed", "hash", hash, "error", err)
		return err
	}

	r.logger.Info("refresh token deleted by hash", "hash", hash)
	return nil
}

func (r *refreshTokenRepository) DeleteByUserID(userID uint) error {
	err := r.db.Unscoped().
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error
	if err != nil {
		r.logger.Error("delete refresh token by user failed", "user_id", userID, "error", err)
		return err
	}

	r.logger.Info("refresh token deleted by user", "user_id", userID)
	return nil
}



