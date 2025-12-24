package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/notification"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"gorm.io/gorm"
)

type AuthService interface {
	RequestPhoneCode(phone string) error
	VerifyPhoneCode(phone string, code string) (userID uint, err error)
}

type authService struct {
	phoneRepo repository.PhoneVerificationRepository
	userRepo  repository.UserRepository
	logger    *slog.Logger
	smsSender notification.SMSSender

	codeTTL     time.Duration
	maxAttempts int
}

func NewAuthService(phoneRepo repository.PhoneVerificationRepository,
	userRepo repository.UserRepository,
	logger *slog.Logger,
	smsSender notification.SMSSender,
) AuthService {
	return &authService{
		phoneRepo:   phoneRepo,
		userRepo:    userRepo,
		smsSender:   smsSender,
		logger:      logger,
		codeTTL:     60 * time.Minute,
		maxAttempts: 20,
	}
}

func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

func generate6Digits() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func compareCodeHash(storedHash, code string) bool {
	return storedHash == hashCode(code)
}

func (s *authService) RequestPhoneCode(phone string) error {
	phone = normalizePhone(phone)

	if phone == "" {
		return errors.New("телефон не может быть пустым")
	}

	code, err := generate6Digits()

	if err != nil {
		s.logger.Error("generate code failed", "phone", phone, "error", err)
		return err
	}

	codeHash := hashCode(code)
	expiresAt := time.Now().Add(s.codeTTL)

	if err := s.phoneRepo.Upsert(phone, codeHash, expiresAt); err != nil {
		s.logger.Error("RequestPhoneCode failed", "phone", phone, "error", err)
		return err
	}

	message := fmt.Sprintf("Ваш код подтверждения: %s", code)

	if err := s.smsSender.Send(phone, message); err != nil {
		s.logger.Error(
			"failed to send sms",
			"phone", phone,
			"error", err,
		)
		return err
	}

	s.logger.Info("sms sent successfully", "phone", phone)

	return nil
}

func (s *authService) VerifyPhoneCode(phone string, code string) (uint, error) {
	phone = normalizePhone(phone)
	code = strings.TrimSpace(code)

	if phone == "" || code == "" {
		return 0, errors.New("phone or code is empty")
	}

	storedHash, expiresAt, attempts, err := s.phoneRepo.GetByPhone(phone)
	if err != nil {
		s.logger.Warn("verification not found", "phone", phone, "error", err)
		return 0, errors.New("code not found")
	}

	if time.Now().After(expiresAt) {
		if err := s.phoneRepo.DeleteByPhone(phone); err != nil {
			s.logger.Error(
				"failed to delete expired verification",
				"phone", phone,
				"error", err,
			)
		}
		return 0, errors.New("code expired")
	}

	if attempts >= s.maxAttempts {
		if err := s.phoneRepo.DeleteByPhone(phone); err != nil {
			s.logger.Error(
				"failed to delete verification after max attempts",
				"phone", phone,
				"attempts", attempts,
				"error", err,
			)
		}

		s.logger.Warn(
			"too many attempts",
			"phone", phone,
			"attempts", attempts,
		)

		return 0, errors.New("too many attempts")
	}

	if !compareCodeHash(storedHash, code) {
		if err := s.phoneRepo.IncrementAttempts(phone); err != nil {
			s.logger.Error(
				"failed to increment attempts",
				"phone", phone,
				"error", err,
			)
		}
		return 0, errors.New("invalid code")
	}

	if err := s.phoneRepo.DeleteByPhone(phone); err != nil {
		s.logger.Error(
			"failed to delete phone verification",
			"phone", phone,
			"error", err,
		)
	}

	user, err := s.userRepo.GetUserByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u := &models.User{Phone: phone}
			if err := s.userRepo.CreateUser(u); err != nil {
				s.logger.Error("create user failed", "phone", phone, "error", err)
				return 0, err
			}
			s.logger.Info("user created on verify", "user_id", u.ID, "phone", phone)
			return u.ID, nil
		}

		s.logger.Error("GetUserByPhone failed", "phone", phone, "error", err)
		return 0, err
	}

	s.logger.Info("phone verified, user found", "user_id", user.ID, "phone", phone)
	return user.ID, nil
}
