package service

import (
	"context"
	"rentalin/internal/errs"
	"rentalin/internal/model"
	"rentalin/internal/repository"
)

type DashboardService interface {
	GetDashboard(ctx context.Context) (*model.Dashboard, error)
}

type dashboardService struct {
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	rentalRepo  repository.RentalRepository
	paymentRepo repository.PaymentRepository
}

func NewDashboardService(userRepo repository.UserRepository, productRepo repository.ProductRepository, rentalRepo repository.RentalRepository, paymentRepo repository.PaymentRepository) DashboardService {
	return &dashboardService{
		userRepo:    userRepo,
		productRepo: productRepo,
		rentalRepo:  rentalRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *dashboardService) GetDashboard(ctx context.Context) (*model.Dashboard, error) {
	totalUsers, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, errs.WrapErr("CountUsers", err)
	}

	totalProducts, err := s.productRepo.Count(ctx)
	if err != nil {
		return nil, errs.WrapErr("CountProducts", err)
	}

	totalRentals, err := s.rentalRepo.Count(ctx)
	if err != nil {
		return nil, errs.WrapErr("CountRentals", err)
	}

	activeRentals, err := s.rentalRepo.CountByStatus(ctx, model.RentalOngoing)
	if err != nil {
		return nil, errs.WrapErr("CountActiveRentals", err)
	}

	pendingPayments, err := s.paymentRepo.CountByStatus(ctx, model.PaymentPending)
	if err != nil {
		return nil, errs.WrapErr("CountPendingPayments", err)
	}

	totalRevenue, err := s.paymentRepo.SumPaidRevenue(ctx)
	if err != nil {
		return nil, errs.WrapErr("SumPaidRevenue", err)
	}

	return &model.Dashboard{
		TotalUsers:      totalUsers,
		TotalProducts:   totalProducts,
		TotalRentals:    totalRentals,
		ActiveRentals:   activeRentals,
		PendingPayments: pendingPayments,
		TotalRevenue:    totalRevenue,
	}, nil
}
