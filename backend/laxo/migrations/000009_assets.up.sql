BEGIN;

CREATE TABLE IF NOT EXISTS assets(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  shop_id CHAR(26) NOT NULL,
  original_filename TEXT,
  extension VARCHAR(32),
  murmur_hash BIGINT,
  file_size INT,
  width INT,
  height INT,
  CONSTRAINT fk_shop_assets FOREIGN KEY(shop_id) REFERENCES shop(id)
);

CREATE TABLE IF NOT EXISTS products_media(
  product_id CHAR(26) NOT NULL,
  asset_id CHAR(26) NOT NULL,
  image_order INT,
  PRIMARY KEY(product_id, asset_id),
  CONSTRAINT fk_products_products_media FOREIGN KEY(product_id) REFERENCES products(id)
  CONSTRAINT fk_assets_products_media FOREIGN KEY(asset_id) REFERENCES assets(id)
);

COMMIT;
