BEGIN;

CREATE TABLE IF NOT EXISTS assets(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  shop_id CHAR(26) NOT NULL,
  murmur_hash CHAR(32) NOT NULL,
  original_filename TEXT,
  extension VARCHAR(32),
  file_size BIGINT,
  width INT,
  height INT,
  UNIQUE(shop_id, murmur_hash),
  CONSTRAINT fk_shop_assets FOREIGN KEY(shop_id) REFERENCES shops(id)
);

CREATE TABLE IF NOT EXISTS products_media(
  product_id CHAR(26) NOT NULL,
  asset_id CHAR(26) NOT NULL,
  image_order INT,
  status VARCHAR(300) DEFAULT 'active' NOT NULL,
  PRIMARY KEY(product_id, asset_id),
  CONSTRAINT fk_products_products_media FOREIGN KEY(product_id) REFERENCES products(id),
  CONSTRAINT fk_assets_products_media FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT status_products_media CHECK(status IN ('active', 'inactive'))
);

COMMIT;
