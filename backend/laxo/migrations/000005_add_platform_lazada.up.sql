BEGIN;
CREATE TABLE IF NOT EXISTS platform_lazada(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  shop_id CHAR(26) NOT NULL,
  access_token VARCHAR(128) NOT NULL,
  country VARCHAR(8) NOT NULL,
  refresh_token VARCHAR(128) NOT NULL,
  account_platform VARCHAR(64) NOT NULL,
  account VARCHAR(300) NOT NULL,
  user_id_vn VARCHAR(48) NOT NULL,
  seller_id_vn VARCHAR(48) NOT NULL,
  short_code_vn VARCHAR(48) NOT NULL,
  refresh_expires_in TIMESTAMP WITH TIME ZONE,
  access_expires_in TIMESTAMP WITH TIME ZONE,
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_shop_platform_lazada FOREIGN KEY(shop_id) REFERENCES shops(id)
);
COMMIT;
