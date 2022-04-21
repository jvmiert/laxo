BEGIN;
ALTER TABLE products_sku_lazada DROP CONSTRAINT fk_sku_products_lazada;
DROP TABLE IF EXISTS products_sku_lazada;

ALTER TABLE products_attribute_lazada DROP CONSTRAINT fk_attribute_products_lazada;
DROP TABLE IF EXISTS products_attribute_lazada;


ALTER TABLE products_lazada DROP CONSTRAINT fk_product_shop_lazada;
DROP TABLE IF EXISTS products_lazada;
COMMIT;
