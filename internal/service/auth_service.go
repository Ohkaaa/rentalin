package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"rentalin/internal/errs"
	"rentalin/internal/model"
	"rentalin/internal/repository"
	"rentalin/pkg/auth"
	"rentalin/pkg/logger"

	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Register(ctx context.Context, username, email, phone, address, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	if strings.HasPrefix(phone, "0") {
		return "62" + phone[1:]
	}
	if strings.HasPrefix(phone, "+62") {
		return "62" + phone[1:]
	}
	return phone
}

func (s *authService) Register(ctx context.Context, username, email, phone, address, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))
	phone = normalizePhone(phone)
	address = strings.TrimSpace(address)
	now := time.Now()

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		logger.Error(
			"failed get user by email",
			err,
			logrus.Fields{
				"email": email,
			},
		)

		return nil, errs.WrapErr("GetByEmail", err)
	}

	if existingUser != nil {
		logger.Warn(
			"user email already exists",
			logrus.Fields{
				"email": email,
			},
		)

		return nil, errs.ErrUserEmailExists
	}

	existingUsername, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		logger.Error(
			"failed get user by username",
			err,
			logrus.Fields{
				"username": username,
			},
		)

		return nil, errs.WrapErr("GetByUsername", err)
	}

	if existingUsername != nil {
		logger.Warn(
			"user username already exists",
			logrus.Fields{
				"username": username,
			},
		)

		return nil, errs.ErrUserUsernameExists
	}

	existingPhone, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		logger.Error(
			"failed get user by phone",
			err,
			logrus.Fields{
				"phone": phone,
			},
		)

		return nil, errs.WrapErr("GetByPhone", err)
	}

	if existingPhone != nil {
		logger.Warn(
			"user phone already exists",
			logrus.Fields{
				"phone": phone,
			},
		)

		return nil, errs.ErrUserPhoneExists
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.Error(
			"failed hash password",
			err,
			nil,
		)

		return nil, errs.WrapErr("HashPassword", err)
	}

	user := &model.User{
		Username:  username,
		Email:     email,
		Phone:     phone,
		Address:   address,
		Password:  hashedPassword,
		Role:      "customer",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error(
			"failed create user",
			err,
			logrus.Fields{
				"email": email,
			},
		)

		return nil, errs.WrapErr("Create", err)
	}

	logger.Info(
		"user registered successfully",
		logrus.Fields{
			"user_id":  user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	)

	user.Password = ""

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		logger.Warn(
			"login failed invalid credentials",
			logrus.Fields{
				"email": email,
			},
		)

		return "", errs.ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"login failed user not found",
				logrus.Fields{
					"email": email,
				},
			)

			return "", errs.ErrInvalidCredentials
		}

		logger.Error(
			"failed get user by email",
			err,
			logrus.Fields{
				"email": email,
			},
		)

		return "", errs.WrapErr("GetByEmail", err)
	}

	if err := auth.CheckPassword(password, user.Password); err != nil {
		logger.Warn(
			"login failed wrong password",
			logrus.Fields{
				"user_id": user.ID,
				"email":   user.Email,
			},
		)

		return "", errs.ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		logger.Error(
			"failed generate jwt",
			err,
			logrus.Fields{
				"user_id": user.ID,
			},
		)

		return "", errs.WrapErr("GenerateJWT", err)
	}

	logger.Info(
		"user login successfully",
		logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
		},
	)

	return token, nil
}
