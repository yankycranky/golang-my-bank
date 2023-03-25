
-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency) 
VALUES ($1,$2,$3) 
returning *;

-- name: GetAccountByName :one
SELECT * FROM accounts
WHERE owner = $1;

-- name: GetAccountByNameForUpdate :one
SELECT * FROM accounts
WHERE owner = $1
FOR NO KEY UPDATE;

-- name: GetAcounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateAccountByOwner :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE owner = sqlc.arg(owner)
returning *;

-- name: DeleteAccountByName :exec
DELETE 
FROM accounts
WHERE owner = $1;

---------------------- TRANSFER -------------------------------


-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount) 
VALUES ($1,$2,$3) 
returning *;

-- name: GetTransferByFromAccount :one
SELECT * FROM transfers
WHERE from_account_id = $1;

-- name: GetTransferByToAccount :one
SELECT * FROM transfers
WHERE to_account_id = $1;

-- name: GetTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateTransferById :exec
UPDATE transfers
SET amount = $2
WHERE id = $1;

-- name: DeleteTransferById :exec
DELETE 
FROM transfers
WHERE id = $1;

---------------------- Entry -------------------------------


-- name: CreateEntry :one
INSERT INTO entries (account_id, amount) 
VALUES ($1,$2) 
returning *;
