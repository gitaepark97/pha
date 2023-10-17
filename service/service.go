package service

import (
	"context"

	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/util"
)

type Service interface {
	// auth
	Register(ctx context.Context, params RegisterParams) (cErr CustomErr)
	Login(ctx context.Context, params LoginParams) (result dto.LoginResponseBody, cErr CustomErr)
	RenewAccessToken(ctx context.Context, params RenewAccessTokenParams) (result dto.RenewAccessTokenResponse, cErr CustomErr)

	// product
	CreateProduct(ctx context.Context, params CreateProductParams) (cErr CustomErr)
	GetProductList(ctx context.Context, params GetProductListParams) (result dto.GetProductListResponse, cErr CustomErr)
	GetProduct(ctx context.Context, params GetProductParams) (result dto.GetProductResponse, cErr CustomErr)
	UpdateProduct(ctx context.Context, params UpdateProductParams) (cErr CustomErr)
	DeleteProduct(ctx context.Context, params DeleteProductParams) (cErr CustomErr)
}

type service struct {
	config     util.Config
	repository repository.Repository
}

func NewService(config util.Config, repository repository.Repository) Service {
	return &service{
		config:     config,
		repository: repository,
	}
}
