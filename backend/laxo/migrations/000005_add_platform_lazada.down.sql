BEGIN;
ALTER TABLE platform_lazada DROP CONSTRAINT fk_shop_platform_lazada;
DROP TABLE IF EXISTS platform_lazada;
COMMIT;
