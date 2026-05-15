-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- -- name: CreateProduct :one
-- INSERT INTO products (
--     name, description, price_in_cents, quantity
-- ) VALUES (
--     $1, $2, $3, $4
-- )
-- RETURNING *;

-- -- name: UpdateProduct :one
-- UPDATE products
-- SET 
--     name = $2,
--     description = $3,
--     price_in_cents = $4,
--     quantity = $5,
--     updated_at = NOW()
-- WHERE id = $1
-- RETURNING *;

-- -- name: DeleteProduct :exec
-- DELETE FROM products
-- WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    customer_id
) VALUES (
    $1
)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, product_id, quantity, price_in_cents
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;