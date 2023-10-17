-- name: CreateProduct :exec
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
);

-- name: GetProductList :many
SELECT
  *
FROM product
WHERE user_id = ?
  AND SearchChosung(name, ?)
ORDER BY created_at DESC
LIMIT 10 OFFSET ?;

-- name: GetProduct :one
SELECT
  *
FROM product
WHERE id = ?;

-- name: UpdateProduct :exec
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
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE
FROM product
WHERE id = ?;