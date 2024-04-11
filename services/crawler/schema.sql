INSERT INTO crawler.products (
    site_code,
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
) SELECT 
    'bomber',
    shop_code,
    product_code,
    name,
    jan,
    price,
    url
FROM bomber_products
ON CONFLICT ON CONSTRAINT products_pkey
DO UPDATE SET
    name = EXCLUDED.name,
    jan = EXCLUDED.jan,
    price = EXCLUDED.price,
    url = EXCLUDED.url;
