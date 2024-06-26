CREATE TABLE IF NOT EXISTS products (
    site_code VARCHAR NOT NULL,
    shop_code VARCHAR NOT NULL,
    product_code VARCHAR NOT NULL,
    name VARCHAR,
    jan VARCHAR,
    price BIGINT,
    url VARCHAR,
    PRIMARY KEY(site_code, shop_code, product_code)
) PARTITION BY LIST (site_code);

CREATE TABLE IF NOT EXISTS products_default PARTITION OF products DEFAULT;
CREATE TABLE IF NOT EXISTS products_ark PARTITION OF products FOR VALUES IN ('ark');
