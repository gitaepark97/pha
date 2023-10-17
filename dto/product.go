package dto

import (
	"time"

	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/util"
)

type CreateProductRequestBody struct {
	Category       string `json:"category" binding:"required,max=100"`
	Price          int32  `json:"price" binding:"required"`
	Cost           int32  `json:"cost" binding:"required"`
	Name           string `json:"name" binding:"required,max=100"`
	Description    string `json:"description" binding:"required"`
	Barcode        string `json:"barcode" binding:"required"`
	ExpirationDate string `json:"expiration_date" binding:"required,date"`
	Size           string `json:"size" binding:"required,product_size"`
}

type GetProductListRequestQuery struct {
	Page    int32  `form:"page" binding:"required,gte=1"`
	Keyword string `form:"keyword" biding:"omitempty"`
}

type GetProductListResponse struct {
	List []GetProductResponse `json:"list"`
}

func NewGetProductListResponse(productList []repository.Product) GetProductListResponse {
	res := GetProductListResponse{}

	for _, product := range productList {
		res.List = append(res.List, NewGetProductResponse(product))
	}

	return res
}

type GetProductRequestPath struct {
	ID int64 `uri:"id" biding:"required"`
}

type GetProductResponse struct {
	ID             int64                  `json:"id"`
	UserID         int64                  `json:"user_id"`
	Category       string                 `json:"category"`
	Price          int32                  `json:"price"`
	Cost           int32                  `json:"cost"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Barcode        string                 `json:"barcode"`
	ExpirationDate string                 `json:"expiration_date"`
	Size           repository.ProductSize `json:"size"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func NewGetProductResponse(product repository.Product) GetProductResponse {
	return GetProductResponse{
		ID:             product.ID,
		UserID:         product.UserID,
		Category:       product.Category,
		Price:          product.Price,
		Cost:           product.Cost,
		Name:           product.Name,
		Description:    product.Description,
		Barcode:        product.Barcode,
		ExpirationDate: product.ExpirationDate.Format(util.DateLayout),
		Size:           product.Size,
		CreatedAt:      product.CreatedAt,
		UpdatedAt:      product.UpdatedAt,
	}
}

type UpdateProductRequestPath = GetProductRequestPath

type UpdateProductRequestBody struct {
	Category       *string `json:"category" binding:"omitempty,max=100"`
	Price          *int32  `json:"price" binding:"omitempty"`
	Cost           *int32  `json:"cost" binding:"omitempty"`
	Name           *string `json:"name" binding:"omitempty,max=100"`
	Description    *string `json:"description" binding:"omitempty"`
	Barcode        *string `json:"barcode" binding:"omitempty"`
	ExpirationDate *string `json:"expiration_date" binding:"omitempty,date"`
	Size           *string `json:"size" binding:"omitempty,product_size"`
}

type DeleteProductRequestPath = GetProductRequestPath
