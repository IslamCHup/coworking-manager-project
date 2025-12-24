package repository

import (
	"log/slog"
	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PhoneVerificationRepository interface {
	Upsert(phone string, codeHash string, expiresAt time.Time) error
	GetByPhone(phone string) (codeHash string, expiresAt time.Time, attempts int, err error)
	IncrementAttempts(phone string) error
	DeleteByPhone(phone string) error
}

type phoneVerificationRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPhoneVerificationRepository(
	db *gorm.DB,
	logger *slog.Logger,
) PhoneVerificationRepository {
	return &phoneVerificationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *phoneVerificationRepository) Upsert(
	phone string,
	codeHash string,
	expiresAt time.Time,
) error {

	verification := models.PhoneVerification{
		Phone:     phone,
		CodeHash:  codeHash,
		ExpiresAt: expiresAt,
		Attempts:  0,
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "phone"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"code_hash":  codeHash,
			"expires_at": expiresAt,
			"attempts":   0,
		}),
	}).Create(&verification).Error

	if err != nil {
		r.logger.Error(
			"PhoneVerification upsert failed",
			"phone", phone,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"PhoneVerification upsert success",
		"phone", phone,
		"expires_at", expiresAt,
	)
	return nil
}

func (r *phoneVerificationRepository) GetByPhone(
	phone string,
) (string, time.Time, int, error) {

	var pv models.PhoneVerification

	err := r.db.
		Where("phone = ?", phone).
		First(&pv).Error

	if err != nil {
		r.logger.Warn(
			"PhoneVerification not found",
			"phone", phone,
			"error", err,
		)
		return "", time.Time{}, 0, err
	}

	r.logger.Info(
		"PhoneVerification fetched",
		"phone", phone,
		"attempts", pv.Attempts,
	)

	return pv.CodeHash, pv.ExpiresAt, pv.Attempts, nil
}

func (r *phoneVerificationRepository) IncrementAttempts(phone string) error {
	err := r.db.Model(&models.PhoneVerification{}).
		Where("phone = ?", phone).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error

	if err != nil {
		r.logger.Error(
			"IncrementAttempts failed",
			"phone", phone,
			"error", err,
		)
		return err
	}

	r.logger.Warn(
		"PhoneVerification attempt incremented",
		"phone", phone,
	)
	return nil
}

func (r *phoneVerificationRepository) DeleteByPhone(phone string) error {
	err := r.db.
		Where("phone = ?", phone).
		Delete(&models.PhoneVerification{}).Error

	if err != nil {
		r.logger.Error(
			"DeletePhoneVerification failed",
			"phone", phone,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"PhoneVerification deleted",
		"phone", phone,
	)
	return nil
}
