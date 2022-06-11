BEGIN;
ALTER TABLE assets DROP CONSTRAINT fk_shop_assets;

ALTER TABLE products_media DROP CONSTRAINT fk_products_products_media;
ALTER TABLE products_media DROP CONSTRAINT fk_assets_products_media;

DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS products_media;
COMMIT;
