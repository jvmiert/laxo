BEGIN;
CREATE TABLE IF NOT EXISTS products(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  name TEXT,
  description TEXT,
  msku TEXT,
  selling_price NUMERIC(21,2),
  cost_price NUMERIC(21,2),
  shop_id CHAR(26) NOT NULL,
  media_id CHAR(26),
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP WITH TIME ZONE,
  UNIQUE (msku, shop_id),
  CONSTRAINT fk_shop_products FOREIGN KEY(shop_id) REFERENCES shops(id)
);

CREATE TABLE IF NOT EXISTS products_media(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  product_id CHAR(26) NOT NULL,
  original_filename TEXT,
  extension VARCHAR(32),
  murmur_hash BIGINT,
  CONSTRAINT fk_products_products_media FOREIGN KEY(product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS products_platform(
  product_id CHAR(26) NOT NULL PRIMARY KEY,
  products_lazada_id CHAR(26),
  CONSTRAINT fk_products_products_platform FOREIGN KEY(product_id) REFERENCES products(id),
  CONSTRAINT fk_products_products_lazada FOREIGN KEY(products_lazada_id) REFERENCES products_lazada(id)
);

ALTER TABLE products ADD CONSTRAINT fk_media_products FOREIGN KEY(media_id) REFERENCES products_media(id);
COMMIT;