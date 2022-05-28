BEGIN;
CREATE TABLE IF NOT EXISTS shops_platforms(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  shop_id CHAR(26) NOT NULL,
  platform_name VARCHAR(300) NOT NULL,
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_shop_platform_shop FOREIGN KEY(shop_id) REFERENCES shops(id),
  CONSTRAINT shops_platforms_platform_names CHECK(platform_name IN ('lazada', 'tiki', 'shopee'))
);
COMMIT;
