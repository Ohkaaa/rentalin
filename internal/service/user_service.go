package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"rentalin/internal/dto"
	"rentalin/internal/errs"
	"rentalin/internal/model"
	"rentalin/internal/repository"
	"rentalin/pkg/logger"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	UpdateProfile(ctx context.Context, userID int, input dto.UpdateUserRequest) error
	DeleteUser(ctx context.Context, userID int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func hidePassword(user *model.User) {
	if user != nil {
		user.Password = ""
	}
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"user profile not found",
				logrus.Fields{
					"user_id": userID,
				},
			)

			return nil, errs.ErrUserNotFound
		}

		logger.Error(
			"failed get user profile",
			err,
			logrus.Fields{
				"user_id": userID,
			},
		)

		return nil, errs.WrapErr("GetByID", err)
	}

	hidePassword(user)

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		logger.Error(
			"failed get all users",
			err,
			nil,
		)

		return nil, errs.WrapErr("GetAll", err)
	}

	for _, user := range users {
		hidePassword(user)
	}

	return users, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, input dto.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"user not found for update",
				logrus.Fields{
					"user_id": userID,
				},
			)

			return errs.ErrUserNotFound
		}

		logger.Error(
			"failed get user before update",
			err,
			logrus.Fields{
				"user_id": userID,
			},
		)

		return errs.WrapErr("GetByID", err)
	}

	if input.Username != nil {
		user.Username = strings.TrimSpace(*input.Username)
	}

	if input.Phone != nil {
		user.Phone = normalizePhone(*input.Phone)
	}

	if input.Address != nil {
		user.Address = strings.TrimSpace(*input.Address)
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error(
			"failed update user profile",
			err,
			logrus.Fields{
				"user_id": userID,
			},
		)

		return errs.WrapErr("Update", err)
	}

	logger.Info(
		"user profile updated",
		logrus.Fields{
			"user_id": userID,
		},
	)

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"user not found for delete",
				logrus.Fields{
					"user_id": userID,
				},
			)

			return errs.ErrUserNotFound
		}

		logger.Error(
			"failed delete user",
			err,
			logrus.Fields{
				"user_id": userID,
			},
		)

		return errs.WrapErr("Delete", err)
	}

	logger.Info(
		"user deleted",
		logrus.Fields{
			"user_id": userID,
		},
	)

	return nil
}
