INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'hikaritv',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM hikaritv_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;

INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'ikebe',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM ikebe_product
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;

INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'kaago',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM kaago_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;
INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'murauchi',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM murauchi_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;


INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'nojima',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM nojima_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;
INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'pc4u',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM pc4u_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;


INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'rakuten',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM rakuten_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;
