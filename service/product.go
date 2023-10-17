package service

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/util"
	"github.com/go-sql-driver/mysql"
)

type CreateProductParams struct {
	UserID int64
	dto.CreateProductRequestBody
}

// 상품 등록 로직
func (service *service) CreateProduct(ctx context.Context, params CreateProductParams) (cErr CustomErr) {
	// string 타입의 날짜 time 타입으로 변환
	parsedTime, err := time.Parse(util.DateLayout, params.ExpirationDate)
	if err != nil {
		cErr = errParseDate
		return
	}

	arg := repository.CreateProductParams{
		UserID:         params.UserID,
		Category:       params.Category,
		Price:          params.Price,
		Cost:           params.Cost,
		Name:           params.Name,
		Description:    params.Description,
		Barcode:        params.Barcode,
		ExpirationDate: parsedTime,
		Size:           repository.ProductSize(params.Size),
	}

	// 상품 생성
	err = service.repository.CreateProduct(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			// 바코드가 중복된 경우
			case repository.DB_DUPLICATE_ERROR:
				switch true {
				case strings.Contains(mysqlErr.Message, "barcode"):
					cErr = errDuplicateBarcode
					return
				}
			// 회원이 없는 경우
			case repository.DB_FK_ERROR:
				switch true {
				case strings.Contains(mysqlErr.Message, "user_id"):
					cErr = errNotFoundUser
					return
				}
			}
		}

		cErr = NewErrInternalServer(err)
		return
	}

	return
}

type GetProductListParams struct {
	UserID int64
	dto.GetProductListRequestQuery
}

// 상품 목록 조회 로직
func (service *service) GetProductList(ctx context.Context, params GetProductListParams) (result dto.GetProductListResponse, cErr CustomErr) {
	arg := repository.GetProductListParams{
		UserID:        params.UserID,
		Searchchosung: params.Keyword,
		Offset:        10 * (params.Page - 1),
	}

	// 상품 검색
	productList, err := service.repository.GetProductList(ctx, arg)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	result = dto.NewGetProductListResponse(productList)
	return
}

type GetProductParams struct {
	UserID int64
	dto.GetProductRequestPath
}

// 상품 상세 조회 로직
func (service *service) GetProduct(ctx context.Context, params GetProductParams) (result dto.GetProductResponse, cErr CustomErr) {
	// 상품 검색
	product, err := service.repository.GetProduct(ctx, params.ID)
	if err != nil {
		// 해당 id의 상품이 없는 경우
		if err == sql.ErrNoRows {
			cErr = errNotFoundProduct
			return
		}

		cErr = NewErrInternalServer(err)
		return
	}

	// 상품 등록 회원 확인
	if product.UserID != params.UserID {
		cErr = errForbiddenProduct
		return
	}

	result = dto.NewGetProductResponse(product)
	return
}

type UpdateProductParams struct {
	UserID int64
	dto.UpdateProductRequestPath
	dto.UpdateProductRequestBody
}

// 상품 수정 로직
func (service *service) UpdateProduct(ctx context.Context, params UpdateProductParams) (cErr CustomErr) {
	// 상품 검색
	product, err := service.repository.GetProduct(ctx, params.ID)
	if err != nil {
		// 해당 id의 상품이 없는 경우
		if err == sql.ErrNoRows {
			cErr = errNotFoundProduct
			return
		}

		cErr = NewErrInternalServer(err)
		return
	}

	// 상품 등록 회원 확인
	if product.UserID != params.UserID {
		cErr = errForbiddenProduct
		return
	}

	arg := repository.UpdateProductParams{
		Category:       product.Category,
		Price:          product.Price,
		Cost:           product.Cost,
		Name:           product.Name,
		Description:    product.Description,
		Barcode:        product.Barcode,
		ExpirationDate: product.ExpirationDate,
		Size:           product.Size,
		ID:             product.ID,
	}

	// mysql의 coalesce 기능 구현
	if params.Category != nil {
		arg.Category = *params.Category
	}
	if params.Price != nil {
		arg.Price = *params.Price
	}
	if params.Name != nil {
		arg.Name = *params.Name
	}
	if params.Description != nil {
		arg.Description = *params.Description
	}
	if params.Barcode != nil {
		arg.Barcode = *params.Barcode
	}
	if params.ExpirationDate != nil {
		// string 타입의 날짜 time 타입으로 변환
		parsedTime, err := time.Parse(util.DateLayout, *params.ExpirationDate)
		if err != nil {
			cErr = errParseDate
			return
		}

		arg.ExpirationDate = parsedTime
	}
	if params.Size != nil {
		arg.Size = repository.ProductSize(*params.Size)
	}

	// 상품 수정
	err = service.repository.UpdateProduct(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			// 바코드가 중복된 경우
			case repository.DB_DUPLICATE_ERROR:
				switch true {
				case strings.Contains(mysqlErr.Message, "barcode"):
					cErr = errDuplicateBarcode
					return
				}
			}
		}

		cErr = NewErrInternalServer(err)
		return
	}

	return
}

type DeleteProductParams struct {
	UserID int64
	dto.DeleteProductRequestPath
}

// 상품 삭제 로직
func (service *service) DeleteProduct(ctx context.Context, params DeleteProductParams) (cErr CustomErr) {
	// 상품 검색
	product, err := service.repository.GetProduct(ctx, params.ID)
	if err != nil {
		// 해당 id의 상품이 없는 경우
		if err == sql.ErrNoRows {
			cErr = errNotFoundProduct
			return
		}

		cErr = NewErrInternalServer(err)
		return
	}

	// 상품 등록 회원 확인
	if product.UserID != params.UserID {
		cErr = errForbiddenProduct
		return
	}

	// 상품 삭제
	err = service.repository.DeleteProduct(ctx, params.ID)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	return
}
