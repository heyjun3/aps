CREATE TABLE IF NOT EXISTS inventories (
    asin varchar NOT NULL,
    fnsku varchar NOT NULL,
    seller_sku varchar PRIMARY KEY,
    condition varchar NOT NULL,
    product_name varchar NOT NULL,
    quantity int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);