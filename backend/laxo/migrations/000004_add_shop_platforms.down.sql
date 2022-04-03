BEGIN;
ALTER TABLE shops_platforms DROP CONSTRAINT fk_shop_platform_shop;
DROP TABLE IF EXISTS shops_platforms;
COMMIT;
