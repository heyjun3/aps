
CREATE TABLE IF NOT EXISTS products_bomber PARTITION OF products FOR VALUES IN ('bomber');
CREATE TABLE IF NOT EXISTS products_hikaritv PARTITION OF products FOR VALUES IN ('hikaritv');
CREATE TABLE IF NOT EXISTS products_ikebe PARTITION OF products FOR VALUES IN ('ikebe');
CREATE TABLE IF NOT EXISTS products_kaago PARTITION OF products FOR VALUES IN ('kaago');
CREATE TABLE IF NOT EXISTS products_murauchi PARTITION OF products FOR VALUES IN ('murauchi');
CREATE TABLE IF NOT EXISTS products_nojima PARTITION OF products FOR VALUES IN ('nojima');
CREATE TABLE IF NOT EXISTS products_pc4u PARTITION OF products FOR VALUES IN ('pc4u');
CREATE TABLE IF NOT EXISTS products_rakuten PARTITION OF products FOR VALUES IN ('rakuten');
