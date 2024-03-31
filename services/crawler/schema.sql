CREATE TABLE products (
    site_code VARCHAR NOT NULL,
    shop_code VARCHAR NOT NULL,
    product_code VARCHAR NOT NULL,
    name VARCHAR,
    jan VARCHAR,
    price BIGINT,
    url VARCHAR,
    PRIMARY KEY(site_code, shop_code, product_code)
)
