BEGIN;
ALTER TABLE orders DROP CONSTRAINT order_platform_name;
ALTER TABLE orders DROP CONSTRAINT fk_social_media_platform_order;
ALTER TABLE orders DROP CONSTRAINT fk_shop_order;

ALTER TABLE order_product DROP CONSTRAINT fk_product_order_product;
ALTER TABLE order_product DROP CONSTRAINT fk_order_order_product;

ALTER TABLE customers DROP CONSTRAINT fk_shop_customer;

ALTER TABLE customers_addresses DROP CONSTRAINT fk_customer_customer_addresses;

ALTER TABLE customers_phones DROP CONSTRAINT fk_customer_customer_phones;

ALTER TABLE customers_socials DROP CONSTRAINT fk_customer_customer_social;
ALTER TABLE customers_socials DROP CONSTRAINT fk_social_platform_customer_social;

DROP TABLE IF EXISTS social_media_platforms;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_product;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS customers_addresses;
DROP TABLE IF EXISTS customers_phones;
DROP TABLE IF EXISTS customers_socials;
COMMIT;
