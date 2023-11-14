CREATE TABLE IF NOT EXISTS lowest_prices (
    seller_sku varchar PRIMARY KEY REFERENCES inventories (seller_sku) ON DELETE CASCADE,
    amount bigint NOT NULL,
    point bigint NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
)