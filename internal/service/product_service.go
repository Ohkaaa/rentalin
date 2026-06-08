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

type ProductService interface {
	CreateProduct(ctx context.Context, name string, dailyPrice int64, stock int) (*model.Product, error)
	GetProductByID(ctx context.Context, productID int) (*model.Product, error)
	GetAllProducts(ctx context.Context) ([]*model.Product, error)
	UpdateProduct(ctx context.Context, productID int, input dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, name string, dailyPrice int64, stock int) (*model.Product, error) {
	name = strings.TrimSpace(name)
	now := time.Now()

	product := &model.Product{
		Name:       name,
		DailyPrice: dailyPrice,
		Stock:      stock,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		logger.Error(
			"failed create product",
			err,
			logrus.Fields{
				"name": product.Name,
			},
		)

		return nil, errs.WrapErr("Create", err)
	}

	logger.Info(
		"product created successfully",
		logrus.Fields{
			"product_id":  product.ID,
			"name":        product.Name,
			"daily_price": product.DailyPrice,
			"stock":       product.Stock,
		},
	)

	return product, nil
}

func (s *productService) GetProductByID(ctx context.Context, productID int) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"product not found",
				logrus.Fields{
					"product_id": productID,
				},
			)

			return nil, errs.ErrProductNotFound
		}

		logger.Error(
			"failed get product",
			err,
			logrus.Fields{
				"product_id": productID,
			},
		)

		return nil, errs.WrapErr("GetByID", err)
	}

	return product, nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		logger.Error(
			"failed get all products",
			err,
			nil,
		)

		return nil, errs.WrapErr("GetAll", err)
	}

	return products, nil
}

func (s *productService) UpdateProduct(ctx context.Context, productID int, input dto.UpdateProductRequest) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"product not found for update",
				logrus.Fields{
					"product_id": productID,
				},
			)

			return errs.ErrProductNotFound
		}

		logger.Error(
			"failed get product before update",
			err,
			logrus.Fields{
				"product_id": productID,
			},
		)

		return errs.WrapErr("GetByID", err)
	}

	if input.Name != nil {
		product.Name = strings.TrimSpace(*input.Name)
	}

	if input.DailyPrice != nil {
		product.DailyPrice = *input.DailyPrice
	}

	if input.Stock != nil {
		product.Stock = *input.Stock
	}

	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		logger.Error(
			"failed update product",
			err,
			logrus.Fields{
				"product_id": productID,
			},
		)

		return errs.WrapErr("Update", err)
	}

	logger.Info(
		"product updated successfully",
		logrus.Fields{
			"product_id":  product.ID,
			"name":        product.Name,
			"daily_price": product.DailyPrice,
			"stock":       product.Stock,
		},
	)

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, productID int) error {
	if err := s.productRepo.Delete(ctx, productID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn(
				"product not found for delete",
				logrus.Fields{
					"product_id": productID,
				},
			)

			return errs.ErrProductNotFound
		}

		logger.Error(
			"failed delete product",
			err,
			logrus.Fields{
				"product_id": productID,
			},
		)

		return errs.WrapErr("Delete", err)
	}

	logger.Info(
		"product deleted successfully",
		logrus.Fields{
			"product_id": productID,
		},
	)

	return nil
}
