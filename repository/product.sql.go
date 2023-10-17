// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: product.sql

package repository

import (
	"context"
	"time"
)

const createProduct = `-- name: CreateProduct :exec
INSERT INTO product(
  user_id,
  category,
  price,
  cost,
  name,
  description,
  barcode,
  expiration_date,
  size
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type CreateProductParams struct {
	UserID         int64       `json:"user_id"`
	Category       string      `json:"category"`
	Price          int32       `json:"price"`
	Cost           int32       `json:"cost"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Barcode        string      `json:"barcode"`
	ExpirationDate time.Time   `json:"expiration_date"`
	Size           ProductSize `json:"size"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error {
	_, err := q.db.ExecContext(ctx, createProduct,
		arg.UserID,
		arg.Category,
		arg.Price,
		arg.Cost,
		arg.Name,
		arg.Description,
		arg.Barcode,
		arg.ExpirationDate,
		arg.Size,
	)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE
FROM product
WHERE id = ?
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT
  id, user_id, category, price, cost, name, description, barcode, expiration_date, size, created_at, updated_at
FROM product
WHERE id = ?
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Category,
		&i.Price,
		&i.Cost,
		&i.Name,
		&i.Description,
		&i.Barcode,
		&i.ExpirationDate,
		&i.Size,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductList = `-- name: GetProductList :many
SELECT
  id, user_id, category, price, cost, name, description, barcode, expiration_date, size, created_at, updated_at
FROM product
WHERE user_id = ?
  AND SearchChosung(name, ?)
ORDER BY created_at DESC
LIMIT 10 OFFSET ?
`

type GetProductListParams struct {
	UserID        int64       `json:"user_id"`
	Searchchosung interface{} `json:"searchchosung"`
	Offset        int32       `json:"offset"`
}

func (q *Queries) GetProductList(ctx context.Context, arg GetProductListParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductList, arg.UserID, arg.Searchchosung, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Category,
			&i.Price,
			&i.Cost,
			&i.Name,
			&i.Description,
			&i.Barcode,
			&i.ExpirationDate,
			&i.Size,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :exec
UPDATE product
SET
  category = ?,
  price = ?,
  cost = ?,
  name = ?,
  description = ?,
  barcode = ?,
  expiration_date = ?,
  size = ?
WHERE id = ?
`

type UpdateProductParams struct {
	Category       string      `json:"category"`
	Price          int32       `json:"price"`
	Cost           int32       `json:"cost"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Barcode        string      `json:"barcode"`
	ExpirationDate time.Time   `json:"expiration_date"`
	Size           ProductSize `json:"size"`
	ID             int64       `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.Category,
		arg.Price,
		arg.Cost,
		arg.Name,
		arg.Description,
		arg.Barcode,
		arg.ExpirationDate,
		arg.Size,
		arg.ID,
	)
	return err
}