BEGIN;
ALTER TABLE products DROP CONSTRAINT fk_shop_products;
ALTER TABLE products DROP CONSTRAINT fk_media_products;

ALTER TABLE products_platform DROP CONSTRAINT fk_products_products_platform;
ALTER TABLE products_platform DROP CONSTRAINT fk_products_products_lazada;

DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS products_platform;
COMMIT;
